# lcma: Legacy code migration assistant

## Legacy Code Modernization Through AI
Fortune 2000 companies face a critical challenge: massive codebases, often decades old, power their core operations. As original developers retire or move on, new engineering teams struggle to decipher these complex legacy systems. This challenge creates significant business risk and slows innovation.

Our project harnesses the power of Generative AI to transform this landscape. By applying AI to legacy code analysis and modernization, we aim to dramatically accelerate the process of understanding, documenting, and converting outdated systems into modern, maintainable codebases built on current technology stacks.

This initiative promises to reduce modernization timelines from years to months, while preserving critical business logic and minimizing risk. Beyond mere code conversion, our AI-driven approach helps capture and transfer the institutional knowledge embedded in these legacy systems, ensuring a smooth transition to modern architectures.


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
- GROQ_API_KEY="YOUR GROQ API KEY"
- MODEL=llama-3.2-90b-vision-preview
- LEGACY_CODE_PATH="YOUR LEGACY CODE PATH DIRECTORY"
- LEGACY_TECH_STACK=[Flask, Python, HTML, CSS, JavaScript]
- MODERN_TECH_STACK=[Golang, Chi, HTMX, Tailwind]
- PROMPT_TEMPLATE_PATH=./prompts
- OUTPUT_FILE_PATH=./legacy_output/output.txt
- REPORT_PATH=./reports
- MODERN_CODE_PATH="YOUR MODERN CODE PATH DIRECTORY"

## FINAL OUTPUT
1. report.md - Gives the full analysis of the legacy code
2. report_code.md - Gives the full code for modern tech stack. 

## BENEFIT
- Significantly reduces the time it takes to convert legacy code to a modern tech stack.
