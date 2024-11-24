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
GROQ_API_KEY="YOUR GROQ API KEY"
MODEL=llama-3.2-90b-vision-preview
LEGACY_CODE_PATH="YOUR LEGACY CODE PATH DIRECTORY"
LEGACY_TECH_STACK=[Flask, Python, HTML, CSS, JavaScript]
MODERN_TECH_STACK=[Golang, Chi, HTMX, Tailwind]
PROMPT_TEMPLATE_PATH=./prompts
OUTPUT_FILE_PATH=./legacy_output/output.txt
REPORT_PATH=./reports
MODERN_CODE_PATH="YOUR MODERN CODE PATH DIRECTORY"

## FINAL OUTPUT
1. report.md - Gives the full analysis of the legacy code
2. report_code.md - Gives the full code for modern tech stack. 

## BENEFIT
- Significantly reduces the time it takes to convert legacy code to a modern tech stack.
