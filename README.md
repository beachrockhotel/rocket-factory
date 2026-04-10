# Rocket-factory

![Coverage](https://img.shields.io/endpoint?url=https://gist.githubusercontent.com/beachrockhotel/e40f85a43679c37c368eaf8e6fe19e97/raw/coverage.json)

Для того чтобы вызывать команды из Taskfile, необходимо установить Taskfile CLI:

```bash
brew install go-task
```

## CI/CD

Проект использует GitHub Actions для непрерывной интеграции и доставки. Основные workflow:

- **CI** (`.github/workflows/ci.yml`) - проверяет код при каждом push и pull request
    - Линтинг кода
    - Проверка безопасности
    - Выполняется автоматическое извлечение версий из Taskfile.yml