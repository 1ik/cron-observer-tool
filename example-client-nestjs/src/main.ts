// IMPORTANT: Load .env file FIRST, before any other imports that might use process.env
import * as dotenv from 'dotenv';
import * as path from 'path';

// Load .env file before creating the app
const envPath = path.resolve(process.cwd(), '.env');
const result = dotenv.config({ path: envPath });

if (result.error) {
  console.warn(`⚠️  Warning: Could not load .env file from ${envPath}`);
  console.warn(`   Error: ${result.error.message}`);
  console.warn(`   Make sure CRON_OBSERVER_API_KEY is set in your environment or .env file`);
} else {
  console.log(`✅ Loaded .env file from ${envPath}`);
}

// Log API key status (without exposing the actual key)
const apiKey = process.env.CRON_OBSERVER_API_KEY;
if (apiKey) {
  console.log(`✅ CRON_OBSERVER_API_KEY is set (length: ${apiKey.length} characters)`);
} else {
  console.error(`❌ CRON_OBSERVER_API_KEY is NOT set!`);
  console.error(`   Please set it in your .env file or as an environment variable`);
}

// Now import other modules after dotenv has loaded
import { NestFactory } from '@nestjs/core';
import { AppModule } from './app.module';

async function bootstrap() {
  const app = await NestFactory.create(AppModule);
  
  // Enable CORS
  app.enableCors();
  
  const port = process.env.PORT || 5202;
  await app.listen(port);
  
  console.log(`🚀 NestJS Example Client running on http://localhost:${port}`);
  console.log(`📡 Execution endpoint: http://localhost:${port}/api/execute`);
}

bootstrap();

