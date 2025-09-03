export default {
  // Общие элементы
  common: {
    loading: 'Загрузка...',
    error: 'Ошибка',
    success: 'Успешно',
    cancel: 'Отмена',
    save: 'Сохранить',
    delete: 'Удалить',
    edit: 'Редактировать',
    close: 'Закрыть',
    confirm: 'Подтвердить',
    back: 'Назад',
    next: 'Далее',
    search: 'Поиск',
    filter: 'Фильтр',
    refresh: 'Обновить',
    export: 'Экспорт',
    import: 'Импорт',
    settings: 'Настройки',
    profile: 'Профиль',
    logout: 'Выход',
    theme: 'Тема',
    language: 'Язык'
  },

  // Навигация
  navigation: {
    home: 'Главная',
    printers: 'Станки',
    reports: 'Отчёты',
    connections: 'Подключения',
    logs: 'Логи'
  },

  // Header
  header: {
    title: 'CNC Manager Pro',
    connected: 'Подключено',
    disconnected: 'Отключено',
    online: 'Онлайн',
    printing: 'В работе',
    uptime: 'Время работы',
    logs: 'Логи',
    systemLogs: 'Системные логи',
    totalLogs: 'Всего',
    updated: 'Обновлено',
    manualUpdate: 'Обновление вручную'
  },

  // Принтеры (ключи оставляем, тексты меняем на CNC)
  printers: {
    title: 'Управление станками',
    subtitle: 'Мониторинг и управление CNC станками',
    status: {
      online: 'Онлайн',
      offline: 'Отключен',
      printing: 'В работе',
      ready: 'Готов',
      error: 'Ошибка'
    },
    details: {
      title: 'Детали станка',
      type: 'Тип',
      version: 'Версия',
      connection: 'Подключение',
      status: 'Статус',
      nozzle: 'Шпиндель',
      bed: 'Стол',
      progress: 'Прогресс',
      timeRemaining: 'Осталось времени'
    },
    controls: {
      printControl: 'Управление задачей',
      gcodeFile: 'G-code файл',
      selectFile: 'Выберите файл',
      startPrint: 'Запустить задачу',
      upload: 'Загрузка...',
      movement: 'Движение по осям',
      step: 'Шаг (мм)',
      quickCommands: 'Быстрые команды',
      customGcode: 'Пользовательский G-code',
      send: 'Отправить',
      sending: 'Отправка команды...'
    },
    commands: {
      homeAll: 'Home All',
      homeXY: 'Home XY',
      homeZ: 'Home Z',
      disableMotors: 'Disable Motors',
      moveAxis: 'Движение по оси',
      commandExecuted: 'Команда выполнена',
      movementCompleted: 'Движение выполнено'
    },
    empty: {
      title: 'Станки не найдены',
      subtitle: 'Подключите станок, чтобы начать работу',
      selectPrinter: 'Выберите станок',
      selectPrinterSubtitle: 'Нажмите на станок в таблице, чтобы просмотреть детальную информацию и управлять им'
    }
  },

  // Подключения
  connections: {
    title: 'Новое подключение',
    subtitle: 'Выберите тип подключения и введите необходимые данные',
    type: 'Тип подключения',
    data: 'Данные подключения',
    ready: 'Готово к подключению',
    connect: 'Подключиться',
    connecting: 'Подключение...',
    help: {
      title: 'Справка по подключению:',
      com: 'COM: Введите номер порта (например, COM3)',
      ip: 'IP: Введите IP-адрес и порт (например, 192.168.1.100:8080)',
      usb: 'USB: Выберите USB-устройство из списка'
    },
    errors: {
      emptyData: 'Пожалуйста, введите данные подключения',
      invalidCom: 'Неверный формат COM порта (например, COM3)',
      invalidIP: 'Неверный формат IP адреса',
      connectionFailed: 'Ошибка подключения',
      success: 'Подключение установлено успешно!'
    }
  },

  // Логи
  logs: {
    title: 'Логи станков',
    search: 'Поиск в логах...',
    filters: {
      all: 'Все',
      info: 'Инфо',
      success: 'Успех',
      warning: 'Предупреждение',
      error: 'Ошибка'
    },
    autoScroll: 'Автопрокрутка',
    noLogs: 'Логи не найдены',
    changeSearch: 'Попробуйте изменить поисковый запрос',
    total: 'Всего',
    shown: 'Показано',
    updateLogs: 'Обновление логов...',
    logsUpdated: 'Логи обновлены',
    logsCleared: 'Логи очищены',
    logsExported: 'Логи экспортированы',
    loadError: 'Ошибка загрузки логов',
    types: {
      system: 'Система',
      connection: 'Подключение',
      printer: 'Станок',
      print: 'Задача',
      command: 'Команда'
    },
    messages: {
      systemStarted: 'Система запущена',
      readyForWork: 'Готов к работе',
      printerConnected: 'Станок подключен успешно',
      printStarted: 'Задача запущена',
      temperatureWarning: 'Температура шпинделя ниже рекомендуемой',
      connectionError: 'Ошибка связи со станком',
      commandExecuted: 'Команда выполнена успешно'
    }
  },

  // Уведомления
  notifications: {
    printerSelected: 'Станок выбран',
    fileRequired: 'Пожалуйста, выберите файл для запуска',
    printerNotSelected: 'Станок не выбран!',
    printStarted: 'Задача запущена успешно!',
    printError: 'Ошибка при отправке запроса на запуск',
    gcodeRequired: 'Введите G-code команду',
    gcodeSent: 'Команда отправлена успешно',
    connectionError: 'Ошибка подключения',
    dataEmpty: 'Данные подключения не могут быть пустыми!'
  },

  // Время
  time: {
    hours: 'ч',
    minutes: 'м',
    seconds: 'с',
    days: 'д'
  },

  // Статусы
  status: {
    loading: 'Загрузка...',
    updating: 'Обновление...',
    connecting: 'Подключение...',
    printing: 'Выполнение...',
    processing: 'Обработка...'
  }
}
