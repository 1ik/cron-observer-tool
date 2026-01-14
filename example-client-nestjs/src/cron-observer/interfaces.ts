export interface ExecutionContext {
  executionId: string;
  taskName: string;
}

export interface CronObserverConfig {
  backendUrl: string;
}

export interface LogEntry {
  message: string;
  level: 'info' | 'warn' | 'error';
}

