# Контекст и рабочая модель проекта oom-watcher

Этот ExecPlan является живым документом. Разделы `Progress`, `Surprises & Discoveries`, `Decision Log` и `Outcomes & Retrospective` должны поддерживаться в актуальном состоянии по мере работы.

Документ нужно поддерживать в соответствии с `.agent/PLANS.md`.

## Purpose / Big Picture

Цель этого документа не в том, чтобы описать новую фичу, а в том, чтобы перенести текущий контекст проекта в один самодостаточный ExecPlan. После чтения этого файла новый участник должен понимать, что именно делает `oom-watcher`, как он запускается, где меняется пользовательское поведение, как собирается `.deb` пакет, как публикуется GitHub Release и как проверить, что приложение работает. Это важно, потому что проект небольшой, но затрагивает сразу несколько слоёв: desktop tray через `systray`, чтение системной памяти из `/proc/meminfo`, конфиг из `/etc`, Debian packaging и release automation через GitHub Actions.

Проверить, что описанный здесь контекст соответствует репозиторию, можно без внешних знаний. Достаточно собрать приложение командой `make build`, упаковать его командой `make release-deb VERSION=<версия>` и затем вручную запустить бинарник в Ubuntu desktop session, где в верхней панели должен появиться процент использования памяти. Если память превышает порог из конфига, этот процент начинает мигать.

## Progress

- [x] (2026-06-02 12:04Z) Собран контекст по структуре репозитория, runtime-логике, конфигу, packaging и release automation.
- [x] (2026-06-02 12:04Z) Создан `.agent/EXECPLAN.md` как единая точка входа для понимания проекта.
- [x] (2026-06-02 12:15Z) Runtime переведен с декоративной иконки на процентный `systray` title; добавлены hot reload конфига и немедленный старт после установки.
- [ ] Поддерживать этот документ в актуальном состоянии при изменении runtime-поведения, формата конфига, packaging или CI release flow.

## Surprises & Discoveries

- Observation: Для показа процентов в верхней панели достаточно `systray.SetTitle` на Linux, поэтому отдельная визуальная иконка больше не нужна.
  Evidence: `github.com/getlantern/systray` имеет `SetTitle` для non-Windows платформ, а `internal/app/systray_app.go` теперь выставляет процентную строку в panel title и мигает именно ей.

- Observation: Для GitHub Release нужен стабильный asset name, иначе установка через `curl` на `latest/download/...` будет ломаться при каждой новой версии.
  Evidence: `Makefile` создает `dist/oom-watcher_amd64.deb`, а `.github/workflows/release.yml` публикует именно этот файл, независимо от внутренней `CalVer`.

- Observation: Debian staging нельзя строить копированием всего каталога `deb`, если внутри него уже лежат ранее собранные `.deb` артефакты или если git не сохранил пустые каталоги.
  Evidence: `Makefile` теперь копирует только `deb/DEBIAN` и `deb/etc`, затем явно создает `usr/local/bin`, чтобы упаковка не зависела от пустого `deb/usr`.

- Observation: В проекте пока нет автоматических тестов.
  Evidence: В репозитории отсутствуют `*_test.go`, а основная верификация сейчас строится на `make build`, `make release-deb` и ручной проверке в desktop session.

## Decision Log

- Decision: Хранить проектный контекст в отдельном файле `.agent/EXECPLAN.md`.
  Rationale: Пользователь попросил перенести контекст проекта в ExecPlan, а правила в `AGENTS.md` и `.agent/PLANS.md` уже задают этот формат как основной для сложного или многослойного контекста.
  Date/Author: 2026-06-02 / Codex

- Decision: Описывать проект как tray application для Ubuntu с ручной runtime-проверкой, а не как headless utility.
  Rationale: Пользовательский эффект виден именно в верхней панели рабочего стола, поэтому контекст должен быть привязан к GUI-поведению, а не только к сборке бинарника.
  Date/Author: 2026-06-02 / Codex

- Decision: Зафиксировать в документе текущий способ установки через `curl -fsSL https://github.com/Syrenny/oom-watcher/releases/latest/download/install.sh | bash`.
  Rationale: Это публичный и user-facing способ доставки, который влияет и на CI, и на release asset naming, и на требования к packaging.
  Date/Author: 2026-06-02 / Codex

- Decision: Явно отметить отсутствие автотестов как часть проектного контекста, а не как временную деталь.
  Rationale: Следующий участник должен понимать, что успешная сборка не равна полной валидации и что ручная проверка tray-поведения остается обязательной.
  Date/Author: 2026-06-02 / Codex

- Decision: Заменить иконку в панели на процентный `systray` title и оставить только прозрачную служебную иконку.
  Rationale: Пользователь явно попросил отказаться от текущей иконки, а текстовый процент проще читается в верхней панели Ubuntu и лучше объясняет смысл приложения.
  Date/Author: 2026-06-02 / Codex

- Decision: Делать hot reload конфига через polling `mtime`, а не через file-system watcher dependency.
  Rationale: Для этого проекта важна надежность и минимализм. Проверка `os.Stat` раз в секунду достаточно проста, не добавляет новую библиотеку и при этом позволяет немедленно подхватывать изменения `poll_interval`, `blink_interval` и порога.
  Date/Author: 2026-06-02 / Codex

- Decision: Запускать приложение сразу после установки из `scripts/install.sh`, а не из Debian `postinst`.
  Rationale: `postinst` выполняется от `root`, а tray-приложение должно запускаться в пользовательской desktop session. Installer script после `sudo dpkg -i` возвращается в пользовательский контекст, где и нужно стартовать `oom-watcher`.
  Date/Author: 2026-06-02 / Codex

## Outcomes & Retrospective

На текущем этапе цель достигнута шире, чем в первой версии документа: кроме переноса контекста в ExecPlan, сам проект уже переведен на более уместную для панели Ubuntu презентацию в виде процентов, умеет автоматически перечитывать конфиг и пытается стартовать сразу после установки. При этом документ все равно не закрывает будущую инженерную работу: в проекте по-прежнему нет автотестов, нет dev-режима для имитации высокой нагрузки и нет отдельной диагностики для desktop environment.

Главный практический вывод из текущего состояния проекта такой: для `oom-watcher` критичны не только Go-исходники, но и целостность цепочки `config -> tray runtime -> Debian package -> GitHub Release assets`. Любое изменение в одной части нужно сверять с соседними слоями.

## Context and Orientation

`oom-watcher` — это минимальное tray-приложение для Ubuntu. Tray application здесь означает маленькую программу без большого окна, которая живет в системной панели сверху и показывает пользователю текущее состояние через panel title, tooltip и пункты выпадающего меню. Конкретно этот проект следит за использованием оперативной памяти. Когда доля занятой памяти становится выше заданного порога, приложение начинает мигать процентным текстом в панели.

Точка входа в приложение находится в `cmd/oom-watcher/main.go`. Этот файл не содержит бизнес-логики. Он только читает конфиг из `/etc/oom-watcher/config.yaml` через `config.NewConfig` и передает его в `internal/app.Run`.

Конфиг описан в `config/config.go`. Сейчас поддерживается только один раздел `memory`, внутри которого есть три поля: `max_used_percent`, `poll_interval` и `blink_interval`. Значения читаются библиотекой `github.com/ilyakaznacheev/cleanenv`, поэтому проект умеет парсить YAML и затем применять env overrides, если они заданы. Валидация встроена в `validate()`: порог должен быть строго между `0` и `100`, интервалы должны быть больше нуля.

GUI-слой начинается в `internal/app/app.go` и `internal/app/systray_app.go`. `app.Run` создает `context.Context`, собирает сервисы и вызывает `systray.Run`. В терминах этого репозитория `systray` — это библиотека `github.com/getlantern/systray`, которая интегрируется с Linux tray/appindicator API и позволяет управлять panel title, tooltip и menu items.

Файл `internal/app/systray_app.go` содержит основную orchestration-логику. `OnReady()` выставляет прозрачную служебную иконку, стартовый текст `--%`, создает отключенный пункт с текущим статусом памяти, отключенный пункт с порогом, отключенный пункт с путем до конфига и отключенный пункт с версией. Затем запускается внутренний цикл `run()`, где живут три `Ticker`: один для опроса памяти, второй для мигания panel title и третий для отслеживания изменения конфига. Если снимок памяти (`service.Snapshot`) говорит, что порог превышен, приложение попеременно показывает и скрывает процентную строку через `systray.SetTitle("")`; если порог не превышен, процент всегда виден.

Отдельных визуальных иконок больше нет. `internal/app/icons.go` генерирует только прозрачную PNG-заглушку, которая нужна `systray` как технический icon payload, а пользовательское значение теперь выводится строкой процента.

Сервисный слой описан в `internal/service/service.go` и `internal/service/gui.go`. Внутри проекта слово “service” означает тонкий слой между tray/UI и источником данных о памяти. `GuiService.UpdateStatus` читает память через `pkg/memory.Read`, обновляет текст пункта меню и возвращает `Snapshot`, где лежат сырые значения (`memory.Stats`), булево поле `Alert` и строка `PanelTitle`, которая потом показывается в верхней панели. `ThresholdTitle` и `Tooltip` формируют пользовательские строки, `SetConfig` подменяет актуальную конфигурацию после reload, а `ShowErr` пишет ошибку в лог и временно заменяет tooltip сообщением об ошибке.

Чтение системной памяти находится в `pkg/memory/reader.go`. Этот модуль открывает `/proc/meminfo`, ищет строки `MemTotal` и `MemAvailable`, затем считает `UsedBytes` как `Total - Available`. Под “used memory” в проекте понимается именно этот расчет, а не значение какого-то отдельного поля ядра. Если одного из ключей нет, `Read()` возвращает ошибку.

Версия бинарника находится в `internal/version/version.go`, а реальное значение подставляется через `-ldflags` из `Makefile`. Это важно для пункта меню `Version: ...` в `internal/app/systray_app.go`.

Packaging для Ubuntu находится в каталоге `deb/`. В `deb/DEBIAN/control` лежат метаданные Debian-пакета и системные зависимости `libayatana-appindicator3-1` и `libgtk-3-0`. В `deb/etc/oom-watcher/config.yaml` лежит конфиг по умолчанию. В `deb/etc/xdg/autostart/oom-watcher.desktop` лежит XDG autostart entry, который указывает на `/usr/local/bin/oom-watcher`. Скрипты `deb/DEBIAN/postinst` и `deb/DEBIAN/postrm` пока минимальны и не содержат сложной логики миграций.

Основная локальная сборка управляется `Makefile`. `make build` форматирует Go-файлы, пытается запустить `golangci-lint`, затем собирает бинарник в `bin/oom-watcher`. `make deb` сначала собирает бинарник, затем создает staging-дерево `.build/deb`, копируя туда только нужные Debian metadata и `etc`, создает `usr/local/bin`, подставляет версию в `DEBIAN/control`, копирует бинарник и вызывает `dpkg-deb --build`. `make release-deb` поверх этого создает каталог `dist/`, куда кладет стабильные release assets `oom-watcher_amd64.deb` и `install.sh`.

Файл `scripts/install.sh` — это публичный installer для GitHub Releases. Он скачивает `https://github.com/Syrenny/oom-watcher/releases/latest/download/oom-watcher_amd64.deb` во временный файл, запускает `sudo dpkg -i`, а если зависимостей не хватает, выполняет `sudo apt-get update` и `sudo apt-get install -f -y`. После успешной установки он пытается остановить старый процесс `oom-watcher` и сразу запустить новый `/usr/local/bin/oom-watcher`, если обнаружена текущая desktop session через `DISPLAY` или `WAYLAND_DISPLAY`.

Automation релиза описан в `.github/workflows/release.yml`. Workflow запускается на `push` в `main` и вручную через `workflow_dispatch`. Он ставит системные build dependencies, вычисляет `CalVer` в формате `YYYY.MM.DD.RUN_NUMBER`, запускает `make release-deb VERSION=<calver>` и публикует два файла через `softprops/action-gh-release@v2`: `dist/oom-watcher_amd64.deb` и `dist/install.sh`. Внешний asset name остается стабильным, а внутренняя версия пакета меняется.

## Plan of Work

Для переноса и поддержки проектного контекста в ExecPlan не нужно отдельно перепридумывать архитектуру. Нужно поддерживать этот файл как “карту проекта”, когда меняется хотя бы один из следующих слоёв: конфиг, panel title поведение, polling/reload логика, Debian package layout, install script, CI release flow. Практически это означает, что после любой заметной правки в `cmd/`, `config/`, `internal/`, `pkg/`, `deb/`, `scripts/` или `.github/workflows/` нужно перечитать затронутые файлы и обновить соответствующие разделы этого документа.

Если в будущем в проект добавятся новые user-facing функции, например новый режим сигнализации, дополнительные пункты меню или телеметрия, сначала нужно описать их здесь на уровне пользовательского эффекта и связать с точными файлами. Если появятся тесты, этот документ нужно дополнить не абстрактным обещанием, а точными командами и именами новых `*_test.go` файлов. Если будет изменен способ доставки, например переход с `curl | bash` на собственный `apt`-репозиторий, разделы `Context and Orientation`, `Concrete Steps`, `Validation and Acceptance` и `Interfaces and Dependencies` должны быть пересмотрены синхронно.

## Concrete Steps

Все команды ниже нужно запускать из корня репозитория `/home/syrenny/Desktop/clones/oom-watcher`.

Для чтения текущего состояния проекта:

    sed -n '1,260p' README.md
    sed -n '1,260p' Makefile
    sed -n '1,260p' .github/workflows/release.yml
    sed -n '1,260p' config/config.go
    sed -n '1,320p' internal/app/systray_app.go
    sed -n '1,320p' internal/service/gui.go
    sed -n '1,260p' pkg/memory/reader.go

Для локальной сборки бинарника:

    make build

Ожидаемый результат — успешная сборка `bin/oom-watcher`. Если `golangci-lint` не установлен, допустимо сообщение:

    golangci-lint not found; skipping lint

Для локальной сборки релизных артефактов с явной версией:

    make release-deb VERSION=2026.06.02.1

Ожидаемый результат — появление двух файлов:

    deb/oom-watcher_2026.06.02.1_amd64.deb
    dist/oom-watcher_amd64.deb
    dist/install.sh

Для ручной установки из latest GitHub Release используется публичная команда:

    curl -fsSL https://github.com/Syrenny/oom-watcher/releases/latest/download/install.sh | bash

Для локального ручного запуска без установки в систему, когда конфиг уже доступен по пути `/etc/oom-watcher/config.yaml`:

    ./bin/oom-watcher

Для ручной проверки hot reload без перезапуска сессии:

    sudoedit /etc/oom-watcher/config.yaml

После сохранения файла приложение должно в течение примерно секунды подхватить новый порог и новые интервалы.

Если нужно проверить содержимое пакета после сборки:

    dpkg-deb -c deb/oom-watcher_2026.06.02.1_amd64.deb

## Validation and Acceptance

Минимальная инженерная валидация проекта сейчас состоит из трех уровней.

Первый уровень — статическая сборка. Нужно выполнить `make build` и убедиться, что бинарник `bin/oom-watcher` создан без ошибок компиляции.

Второй уровень — packaging и release assets. Нужно выполнить `make release-deb VERSION=<версия>` и убедиться, что:

1. создался versioned `.deb` в `deb/`;
2. создался стабильный `dist/oom-watcher_amd64.deb`;
3. создался `dist/install.sh`;
4. `deb/DEBIAN/control` в рабочем дереве не был переписан напрямую, потому что version injection делается внутри `.build/deb`.

Третий уровень — пользовательское поведение в Ubuntu desktop session. После установки пакета или ручного запуска бинарника в верхней панели должен появиться процент использования памяти. В меню должны быть видны четыре неактивных пункта: текущая память, порог, путь до конфига и версия. Если временно понизить `memory.max_used_percent` в `/etc/oom-watcher/config.yaml` до значения ниже текущей фактической загрузки памяти, процент должен начать мигать без перезапуска сессии. Если изменить `poll_interval` или `blink_interval`, приложение должно подхватить их автоматически примерно за секунду. Если вернуть порог выше фактической загрузки, мигание должно прекратиться, а процент снова должен отображаться постоянно.

Пока в репозитории нет unit tests или integration tests, нельзя считать задачу полностью проверенной только по `make build`. Для любых изменений, затрагивающих `internal/app/systray_app.go`, `internal/app/icons.go`, `internal/service/gui.go`, `scripts/install.sh` или `pkg/memory/reader.go`, ручная desktop-проверка обязательна.

## Idempotence and Recovery

Команды `make build` и `make release-deb VERSION=<версия>` спроектированы как повторяемые. Они перезаписывают локальные build artifacts в `bin/`, `.build/`, `dist/` и versioned `.deb` в `deb/`. Если упаковка падает на половине пути, безопасный способ повторить процесс — снова вызвать ту же команду; `make deb` и `make release-deb` начинают с очистки staging-каталога `.build/deb`.

Если ручная установка через `install.sh` завершилась неудачно на этапе `dpkg -i`, скрипт автоматически пытается исправить зависимости через `apt-get install -f -y`. Если и после этого пакет не установился, нужно отдельно повторить:

    sudo dpkg -i /tmp/<скачанный_файл>.deb
    sudo apt-get install -f -y

и уже по сообщению `dpkg` разбирать конкретную системную проблему. Если пакет установился, но процент в панели не появился сразу, нужно проверить, была ли у shell-сессии переменная `DISPLAY` или `WAYLAND_DISPLAY`; при их отсутствии installer корректно пропускает немедленный запуск и оставляет только autostart на следующий login. Этот документ не предполагает destructive rollback steps, потому что в проекте пока нет миграций данных или long-lived state.

## Artifacts and Notes

Ключевые артефакты, которые подтверждают текущее устройство проекта:

    README.md
    Makefile
    .github/workflows/release.yml
    scripts/install.sh
    config/config.go
    internal/app/systray_app.go
    internal/app/icons.go
    internal/service/gui.go
    pkg/memory/reader.go
    deb/DEBIAN/control
    deb/etc/oom-watcher/config.yaml
    deb/etc/xdg/autostart/oom-watcher.desktop

Короткий пример ожидаемого вывода после успешной локальной сборки релизных артефактов:

    Built package: deb/oom-watcher_2026.06.02.1_amd64.deb
    Built release asset: dist/oom-watcher_amd64.deb
    Built installer: dist/install.sh

## Interfaces and Dependencies

Проект написан на Go 1.25 и использует две прямые зависимости.

Первая — `github.com/getlantern/systray v1.2.2`. Она отвечает за интеграцию с tray/appindicator и используется в `internal/app/app.go`, `internal/app/systray_app.go` и `internal/service/gui.go`. Ключевые вызовы, которые должны существовать и оставаться рабочими: `systray.Run`, `systray.SetIcon`, `systray.SetTitle`, `systray.SetTooltip`, `systray.AddMenuItem`, `(*systray.MenuItem).SetTitle` и `(*systray.MenuItem).Disable`.

Вторая — `github.com/ilyakaznacheev/cleanenv v1.5.0`. Она отвечает за загрузку конфигурации и используется в `config/config.go`. Ключевой интерфейс проекта на этом слое — функция `config.NewConfig(configPath string) (*Config, error)`, которая читает YAML и применяет env overrides.

Внутренние интерфейсы и типы, на которые нужно ориентироваться при изменениях:

В `internal/service/service.go` объявлен интерфейс:

    type Gui interface {
        UpdateStatus(statusItem *systray.MenuItem) (Snapshot, error)
        ThresholdTitle() string
        Tooltip(snapshot Snapshot) string
        SetConfig(cfg config.Config)
        ShowErr(err error)
    }

В `internal/service/gui.go` определен снимок пользовательского состояния:

    type Snapshot struct {
        Stats      memory.Stats
        Alert      bool
        PanelTitle string
    }

В `pkg/memory/reader.go` определен тип данных чтения памяти:

    type Stats struct {
        TotalBytes     uint64
        AvailableBytes uint64
        UsedBytes      uint64
        UsedPercent    float64
    }

В `internal/app/systray_app.go` основной цикл полагается на то, что `Gui.UpdateStatus` возвращает полный `Snapshot`, а не только строку статуса. Это важно, потому что решение о мигании, показе panel title и применении нового конфига принимается на уровне app orchestration, а не внутри сервиса.

Из системных зависимостей для Ubuntu важны `libayatana-appindicator3-1` и `libgtk-3-0` при установке пакета, а для локальной сборки — `build-essential`, `pkg-config`, `libgtk-3-dev` и `libayatana-appindicator3-dev`. Эти пакеты уже отражены соответственно в `deb/DEBIAN/control`, `README.md` и `.github/workflows/release.yml`.

Изменение в этом документе от 2026-06-02: документ обновлен под новую модель интерфейса и runtime, где верхняя панель показывает проценты вместо декоративной иконки, конфиг перечитывается автоматически, а installer пытается запускать приложение сразу после установки.
