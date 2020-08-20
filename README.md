# Dockle Run Application

## Parameters

| Environment Variables    | Description                                                           | Required | Default | Example                 |
| ------------------------ | --------------------------------------------------------------------- | -------- | ------- | ----------------------- |
| SCAN_IMAGE               | Docker Images to scan                                                 | \*       |         |                         |
| IMAGE_SCAN_TYPE          | Types to scan [ecr, dockle]                                           | \*       |         |                         |
| IS_LOCAL_IMAGE           | A locally existing Image or                                           |          | false   |                         |
| EXIT_DOCKLE_ERROR_LEVEL  | Dockle error level to return an error [fatal, warn, info, skip, pass] |          | fatal   |                         |
| ISSUE_DOCKLE_ERROR_LEVEL | Dockle error level to create an issue [fatal, warn, info, skip, pass] |          | warn    |                         |
| IGNORE_ERROR_CODES       | Error codes to ignore                                                 |          | []      | CIS-DI-0000:DKL-DI-0000 |
| GIT_HUB_ACCESS_TOKEN     | GitHub access tokens                                                  | \*       |         |                         |
| GIT_HUB_Owner            | Name of the owner or organization of the repository                   | \*       |         |                         |
| GIT_HUB_Repository       | Repository name                                                       | \*       |         |                         |
| ISSUE_APPLICATION_TYPE   | The application name to be set for the label                          | \*       |         |                         |
| ISSUE_ENVIRONMENT        | The environment to be labeled                                         | \*       |         |                         |
