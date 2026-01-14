import { Injectable, Logger } from '@nestjs/common';
import { LogEntry } from './interfaces';

@Injectable()
export class CronObserverService {
  private readonly logger = new Logger(CronObserverService.name);
  private backendUrl: string;

  constructor(backendUrl: string) {
    this.backendUrl = backendUrl;
  }

  async log(executionId: string, message: string, level: 'info' | 'warn' | 'error' = 'info'): Promise<void> {
    const logEntry: LogEntry = {
      message,
      level,
    };

    try {
      const url = `${this.backendUrl}/api/v1/executions/${executionId}/logs`;
      const response = await fetch(url, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
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
      const url = `${this.backendUrl}/api/v1/executions/${executionId}/status`;
      this.logger.log(`Attempting to mark execution ${executionId} as SUCCESS`);
      this.logger.log(`Sending PATCH request to: ${url}`);
      
      const response = await fetch(url, {
        method: 'PATCH',
        headers: {
          'Content-Type': 'application/json',
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

