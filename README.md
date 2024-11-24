# lcma: Legacy code migration assistant

## Code Migration Assistant
A code migration tool that analyzes codebase and documentation of a legacy Flask/Python application and suggests a plan to systematically migrate it to a cloud native Golang Chi HTMX Tailwind based responsive web application.

Features:
- Use layered application design that uses responsive design, horizontally scaleable & uses concurrency for performance.
- Uses modern security, DB access, configuration, decoupled design & error handling patterns
- Automatic identification of repetitive code patterns in Flask/Python codebase that could benefit from redesign/refactoring
- Smart recommendations for API integration points
- Risk assessment for each suggested modification
- Generates boilerplate code for Golang Chi HTMX Tailwind based responsive web application

## User would need to configure these configuration input in the .env file:
GROQ_API_KEY=
MODEL=llama-3.1-70b-versatile
LEGACY_CODE_PATH=
LEGACY_TECH_STACK=[]
MODERN_TECH_STACK=[]
PROMPT_TEMPLATE_PATH=./prompt.txt
OUTPUT_PATH=./output.txt
REPORT_PATH=./report.md
MODERN_CODE_PATH=