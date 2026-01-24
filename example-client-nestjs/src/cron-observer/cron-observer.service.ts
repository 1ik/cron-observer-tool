import { Injectable, Logger } from '@nestjs/common';
import { LogEntry } from './interfaces';

@Injectable()
export class CronObserverService {
  private readonly logger = new Logger(CronObserverService.name);
  private backendUrl: string;
  private apiKey: string;

  constructor(backendUrl: string, apiKey: string) {
    this.backendUrl = backendUrl;
    this.apiKey = apiKey;
    if (!apiKey) {
      this.logger.warn('⚠️  API key is empty! Set CRON_OBSERVER_API_KEY environment variable.');
    }
  }

  async log(executionId: string, message: string, level: 'info' | 'warn' | 'error' = 'info'): Promise<void> {
    const logEntry: LogEntry = {
      message,
      level,
    };

    try {
      if (!this.apiKey) {
        this.logger.error(`Cannot send log: API key is not set. Set CRON_OBSERVER_API_KEY environment variable.`);
        return;
      }

      const url = `${this.backendUrl}/api/v1/executions/${executionId}/logs`;
      const response = await fetch(url, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': this.apiKey,
        },
        body: JSON.stringify(logEntry),
      });

      if (!response.ok) {
        const errorText = await response.text();
        this.logger.error(`Failed to send log to cron-observer: ${response.status} - ${errorText}`);
      }
    } catch (error) {
      this.logger.error(`Error sending log to cron-observer: ${error.message}`, error.stack);
    }
  }

  /**
   * Mark execution as successful
   * @param executionId - The execution UUID
   */
  async success(executionId: string): Promise<void> {
    try {
      if (!this.apiKey) {
        this.logger.error(`Cannot mark execution as SUCCESS: API key is not set. Set CRON_OBSERVER_API_KEY environment variable.`);
        throw new Error('API key is not set');
      }

      const url = `${this.backendUrl}/api/v1/executions/${executionId}/status`;
      this.logger.log(`Attempting to mark execution ${executionId} as SUCCESS`);
      this.logger.log(`Sending PATCH request to: ${url}`);
      
      const response = await fetch(url, {
        method: 'PATCH',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': this.apiKey,
        },
        body: JSON.stringify({ status: 'SUCCESS' }),
      });

      const responseText = await response.text();
      
      if (!response.ok) {
        this.logger.error(`Failed to update execution status to SUCCESS: ${response.status} - ${responseText}`);
        throw new Error(`HTTP ${response.status}: ${responseText}`);
      } else {
        this.logger.log(`✅ Execution ${executionId} marked as SUCCESS`);
        this.logger.log(`Response: ${responseText}`);
      }
    } catch (error) {
      this.logger.error(`❌ Error updating execution status: ${error.message}`, error.stack);
      throw error; // Re-throw so caller knows it failed
    }
  }

  /**
   * Mark execution as failed
   * @param executionId - The execution UUID
   * @param errorMessage - Optional error message
   */
  async fail(executionId: string, errorMessage?: string): Promise<void> {
    try {
      if (!this.apiKey) {
        this.logger.error(`Cannot mark execution as FAILED: API key is not set. Set CRON_OBSERVER_API_KEY environment variable.`);
        throw new Error('API key is not set');
      }

      const url = `${this.backendUrl}/api/v1/executions/${executionId}/status`;
      const body: { status: string; error?: string } = { status: 'FAILED' };
      if (errorMessage) {
        body.error = errorMessage;
      }

      this.logger.log(`Attempting to mark execution ${executionId} as FAILED`);
      this.logger.log(`Sending PATCH request to: ${url}`);
      
      const response = await fetch(url, {
        method: 'PATCH',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': this.apiKey,
        },
        body: JSON.stringify(body),
      });

      const responseText = await response.text();

      if (!response.ok) {
        this.logger.error(`Failed to update execution status to FAILED: ${response.status} - ${responseText}`);
        throw new Error(`HTTP ${response.status}: ${responseText}`);
      } else {
        this.logger.log(`❌ Execution ${executionId} marked as FAILED${errorMessage ? `: ${errorMessage}` : ''}`);
        this.logger.log(`Response: ${responseText}`);
      }
    } catch (error) {
      this.logger.error(`Error updating execution status: ${error.message}`, error.stack);
      throw error; // Re-throw so caller knows it failed
    }
  }
}

