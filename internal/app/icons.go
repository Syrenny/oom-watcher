package app

/*
#cgo pkg-config: pangocairo

#include <stdlib.h>
#include <string.h>
#include <cairo.h>
#include <pango/pangocairo.h>

typedef struct {
	unsigned char *data;
	size_t len;
	size_t cap;
} png_buffer_t;

static cairo_status_t write_png_chunk(void *closure, const unsigned char *data, unsigned int length) {
	png_buffer_t *buf = (png_buffer_t *)closure;
	size_t needed = buf->len + length;
	if (needed > buf->cap) {
		size_t new_cap = buf->cap == 0 ? needed : buf->cap * 2;
		if (new_cap < needed) {
			new_cap = needed;
		}
		unsigned char *new_data = (unsigned char *)realloc(buf->data, new_cap);
		if (new_data == NULL) {
			return CAIRO_STATUS_NO_MEMORY;
		}
		buf->data = new_data;
		buf->cap = new_cap;
	}
	memcpy(buf->data + buf->len, data, length);
	buf->len += length;
	return CAIRO_STATUS_SUCCESS;
}

static int render_text_png(const char *text, int visible, unsigned char **out_data, int *out_len, char **out_err) {
	const int pad_x = 3;
	const int pad_y = 2;
	const char *font_desc_text = "Ubuntu Sans 11";

	cairo_surface_t *measure_surface = cairo_image_surface_create(CAIRO_FORMAT_ARGB32, 1, 1);
	cairo_t *measure_cr = cairo_create(measure_surface);
	PangoLayout *measure_layout = pango_cairo_create_layout(measure_cr);
	PangoFontDescription *font_desc = pango_font_description_from_string(font_desc_text);
	pango_layout_set_font_description(measure_layout, font_desc);
	pango_layout_set_text(measure_layout, text, -1);

	int text_w = 0;
	int text_h = 0;
	pango_layout_get_pixel_size(measure_layout, &text_w, &text_h);

	g_object_unref(measure_layout);
	cairo_destroy(measure_cr);
	cairo_surface_destroy(measure_surface);

	int width = text_w + pad_x * 2;
	int height = text_h + pad_y * 2;
	if (width < 1) {
		width = 1;
	}
	if (height < 1) {
		height = 1;
	}

	cairo_surface_t *surface = cairo_image_surface_create(CAIRO_FORMAT_ARGB32, width, height);
	if (cairo_surface_status(surface) != CAIRO_STATUS_SUCCESS) {
		pango_font_description_free(font_desc);
		*out_err = strdup("failed to create cairo surface");
		return 1;
	}

	cairo_t *cr = cairo_create(surface);
	cairo_set_operator(cr, CAIRO_OPERATOR_SOURCE);
	cairo_set_source_rgba(cr, 0, 0, 0, 0);
	cairo_paint(cr);

	if (visible) {
		PangoLayout *layout = pango_cairo_create_layout(cr);
		pango_layout_set_font_description(layout, font_desc);
		pango_layout_set_text(layout, text, -1);
		cairo_move_to(cr, pad_x, pad_y);
		cairo_set_source_rgba(cr, 0.96, 0.97, 0.98, 1.0);
		pango_cairo_show_layout(cr, layout);
		g_object_unref(layout);
	}

	pango_font_description_free(font_desc);

	png_buffer_t buf = {0};
	cairo_status_t status = cairo_surface_write_to_png_stream(surface, write_png_chunk, &buf);
	cairo_destroy(cr);
	cairo_surface_destroy(surface);

	if (status != CAIRO_STATUS_SUCCESS) {
		free(buf.data);
		*out_err = strdup(cairo_status_to_string(status));
		return 1;
	}

	*out_data = buf.data;
	*out_len = (int)buf.len;
	return 0;
}
*/
import "C"

import (
	"fmt"
	"unsafe"
)

type panelIcons struct {
	Visible []byte
	Blank   []byte
}

func renderPanelIcons(text string) (panelIcons, error) {
	visible, err := renderPanelIcon(text, true)
	if err != nil {
		return panelIcons{}, err
	}

	blank, err := renderPanelIcon(text, false)
	if err != nil {
		return panelIcons{}, err
	}

	return panelIcons{
		Visible: visible,
		Blank:   blank,
	}, nil
}

func renderPanelIcon(text string, visible bool) ([]byte, error) {
	visibleFlag := C.int(0)
	if visible {
		visibleFlag = 1
	}

	cText := C.CString(text)
	defer C.free(unsafe.Pointer(cText))

	var cData *C.uchar
	var cLen C.int
	var cErr *C.char
	if C.render_text_png(cText, visibleFlag, &cData, &cLen, &cErr) != 0 {
		defer C.free(unsafe.Pointer(cErr))
		return nil, fmt.Errorf("render panel icon: %s", C.GoString(cErr))
	}
	defer C.free(unsafe.Pointer(cData))

	return C.GoBytes(unsafe.Pointer(cData), cLen), nil
}
