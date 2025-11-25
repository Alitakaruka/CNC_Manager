export default {
  // Common elements
  common: {
    loading: 'Loading...',
    error: 'Error',
    success: 'Success',
    cancel: 'Cancel',
    save: 'Save',
    delete: 'Delete',
    edit: 'Edit',
    close: 'Close',
    confirm: 'Confirm',
    back: 'Back',
    next: 'Next',
    search: 'Search',
    filter: 'Filter',
    refresh: 'Refresh',
    export: 'Export',
    import: 'Import',
    settings: 'Settings',
    profile: 'Profile',
    logout: 'Logout',
    theme: 'Theme',
    language: 'Language'
  },

  // Navigation
  navigation: {
    home: 'Home',
    printers: 'Machines',
    reports: 'Reports',
    connections: 'Connections',
    logs: 'Logs'
  },

  // Header
  header: {
    title: 'CNC Manager Pro',
    connected: 'Connected',
    disconnected: 'Disconnected',
    online: 'Online',
    printing: 'Running',
    uptime: 'Uptime',
    logs: 'Logs',
    systemLogs: 'System Logs',
    totalLogs: 'Total',
    updated: 'Updated',
    manualUpdate: 'Manual Update'
  },

  // Printers (keys kept for compatibility)
  printers: {
    title: 'Machine Management',
    subtitle: 'Monitor and manage CNC machines',
    status: {
      online: 'Online',
      offline: 'Offline',
      printing: 'Running',
      ready: 'Ready',
      error: 'Error'
    },
    details: {
      title: 'Machine Details',
      type: 'Type',
      version: 'Version',
      connection: 'Connection',
      status: 'Status',
      nozzle: 'Spindle',
      bed: 'Table',
      progress: 'Progress',
      timeRemaining: 'Time Remaining',
      reconnect: 'Reconnect'
    },
    controls: {
      printControl: 'Job Control',
      gcodeFile: 'G-code File',
      selectFile: 'Select File',
      startPrint: 'Start Job',
      upload: 'Uploading...',
      movement: 'Axis Movement',
      step: 'Step (mm)',
      quickCommands: 'Quick Commands',
      customGcode: 'Custom G-code',
      send: 'Send',
      sending: 'Sending Command...'
    },
    commands: {
      homeAll: 'Home All',
      homeXY: 'Home XY',
      homeZ: 'Home Z',
      disableMotors: 'Disable Motors',
      moveAxis: 'Move Axis',
      commandExecuted: 'Command Executed',
      movementCompleted: 'Movement Completed'
    },
    empty: {
      title: 'No Machines Found',
      subtitle: 'Connect a machine to get started',
      selectPrinter: 'Select Machine',
      selectPrinterSubtitle: 'Click a machine in the table to view details and control it'
    }
  },

  // Connections
  connections: {
    title: 'New Connection',
    subtitle: 'Select connection type and enter required data',
    type: 'Connection Type',
    data: 'Connection Data',
    ready: 'Ready to Connect',
    connect: 'Connect',
    connecting: 'Connecting...',
    help: {
      title: 'Connection Help:',
      com: 'COM: Enter port number (e.g., COM3)',
      ip: 'IP: Enter IP address and port (e.g., 192.168.1.100:8080)',
      usb: 'USB: Select USB device from list'
    },
    errors: {
      emptyData: 'Please enter connection data',
      invalidCom: 'Invalid COM port format (e.g., COM3)',
      invalidIP: 'Invalid IP address format',
      connectionFailed: 'Connection failed',
      success: 'Connection established successfully!'
    }
  },

  // Logs
  logs: {
    title: 'Machine Logs',
    search: 'Search in logs...',
    filters: {
      all: 'All',
      info: 'Info',
      success: 'Success',
      warning: 'Warning',
      error: 'Error'
    },
    autoScroll: 'Auto-scroll',
    noLogs: 'No logs found',
    changeSearch: 'Try changing your search query',
    total: 'Total',
    shown: 'Shown',
    updateLogs: 'Updating logs...',
    logsUpdated: 'Logs updated',
    logsCleared: 'Logs cleared',
    logsExported: 'Logs exported',
    loadError: 'Error loading logs',
    types: {
      system: 'System',
      connection: 'Connection',
      printer: 'Machine',
      print: 'Job',
      command: 'Command'
    },
    messages: {
      systemStarted: 'System started',
      readyForWork: 'Ready for work',
      printerConnected: 'Machine connected successfully',
      printStarted: 'Job started',
      temperatureWarning: 'Spindle temperature below recommended',
      connectionError: 'Machine connection error',
      commandExecuted: 'Command executed successfully'
    }
  },

  // Notifications
  notifications: {
    printerSelected: 'Machine selected',
    fileRequired: 'Please select a file to start',
    printerNotSelected: 'Machine not selected!',
    printStarted: 'Job started successfully!',
    printError: 'Error sending start request',
    gcodeRequired: 'Enter G-code command',
    gcodeSent: 'Command sent successfully',
    connectionError: 'Connection error',
    dataEmpty: 'Connection data cannot be empty!'
  },

  // Time
  time: {
    hours: 'h',
    minutes: 'm',
    seconds: 's',
    days: 'd'
  },

  // Status
  status: {
    loading: 'Loading...',
    updating: 'Updating...',
    connecting: 'Connecting...',
    printing: 'Running...',
    processing: 'Processing...'
  }
}
