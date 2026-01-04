# @cron-observer/lib

Shared utilities and API client for the Cron Observer frontend.

## API Client

The API client is automatically generated from the backend's OpenAPI specification using [openapi-zod-client](https://github.com/astahmer/openapi-zod-client).

### Generating the API Client

To regenerate the API client after backend API changes:

```bash
pnpm gen:api
```

This will:
1. Convert the Swagger 2.0 spec from `backend/api-docs/swagger.json` to OpenAPI 3.0 format
2. Save the OpenAPI 3.0 spec to `backend/api-docs/openapi.json`
3. Generate a typed Zodios client with proper request/response types in `src/api-client.ts`

### Using the API Client

```typescript
import { createApiClient } from '@cron-observer/lib';

// Create an API client instance
const api = createApiClient('http://localhost:8080/api/v1');

// Use typed API methods
const projects = await api.getProjects();
const newProject = await api.postProjects({
  name: 'My Project',
  description: 'Project description'
});
```

### OpenAPI 3.0 Conversion

The backend generates a Swagger 2.0 spec, which is automatically converted to OpenAPI 3.0 during the generation process using `swagger2openapi`. This ensures full compatibility with `openapi-zod-client` and proper generation of request body types, response schemas, and Zod validation.

