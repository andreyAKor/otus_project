# Проектная работа по курсу "Разработчик Golang"
[![Build Status](https://travis-ci.com/andreyAKor/otus_project.svg?branch=master)](https://travis-ci.com/andreyAKor/otus_project)

Все команды выполняются из корня проекта.

### Запуск сервиса
Поднимаем docker-compose с сервисом внутри
```shell script
$ make run
docker-compose up --build
Creating network "otus_project_default" with the default driver
Building image_previewer
...
Successfully built cbc5e0e4e6a1
Successfully tagged otus_project_image_previewer:latest
Creating otus_project_image_previewer_1 ... done
Attaching to otus_project_image_previewer_1
```

Сервис поднимается на локальном хосте на порту 6080.
По умолчанию настроено кеширование только 3-х последних нарезанных изображений. Все нарезанные изображения сохраняются в папке `/tmp`, а имена закешированных файлов имеют префикс `image_previewer` (пример: `image_previewer564429678`).

Имеется исходное изображение http://www.audubon.org/sites/default/files/a1_1902_16_barred-owl_sandra_rothenberg_kk.jpg размером 1714px × 1207px.

В браузере проверяем различные варианты нарезки этого изображение:
- размер нарезки 428px × 301px: http://localhost:6080/428/301/www.audubon.org/sites/default/files/a1_1902_16_barred-owl_sandra_rothenberg_kk.jpg
- размер нарезки 500px × 200px: http://localhost:6080/500/200/www.audubon.org/sites/default/files/a1_1902_16_barred-owl_sandra_rothenberg_kk.jpg
- размер нарезки 300px × 200px: http://localhost:6080/300/200/www.audubon.org/sites/default/files/a1_1902_16_barred-owl_sandra_rothenberg_kk.jpg
- размер нарезки 200px × 500px: http://localhost:6080/200/500/www.audubon.org/sites/default/files/a1_1902_16_barred-owl_sandra_rothenberg_kk.jpg
- размер нарезки 2000px × 200px: http://localhost:6080/2000/200/www.audubon.org/sites/default/files/a1_1902_16_barred-owl_sandra_rothenberg_kk.jpg
- размер нарезки 2000px × 1408px: http://localhost:6080/2000/1408/www.audubon.org/sites/default/files/a1_1902_16_barred-owl_sandra_rothenberg_kk.jpg
- размер нарезки 3428px × 2414px: http://localhost:6080/3428/2414/www.audubon.org/sites/default/files/a1_1902_16_barred-owl_sandra_rothenberg_kk.jpg

### Прочие операции с make
- `make install` - устанавливает все необходимые модули через go mod
- `make generate` - go-генерация небходимых для проекта пакетов
- `make lint` - прогонка проекта линтером
- `make build` - сборка сервиса
- `make run` - запуск сервсиа в docker-контейнере через docker-compose
- `make run-dev` - сборка и запуск сервиса для нужд разработкиб без использования docker-контейнера
- `make test` - запуст юнит-тестов
- `make test-integration` - запуст интеграционных тестов
