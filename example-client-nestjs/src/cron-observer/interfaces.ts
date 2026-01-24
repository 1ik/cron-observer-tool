export interface ExecutionContext {
  executionId: string;
  taskName: string;
}

export interface CronObserverConfig {
  backendUrl: string;
  apiKey: string;
}

export interface LogEntry {
  message: string;
  level: 'info' | 'warn' | 'error';
}

