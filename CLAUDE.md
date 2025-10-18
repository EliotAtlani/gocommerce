# Claude Instructions for GoCommerce Project

## Project Context
This is a **learning project** for mastering Go through building a microservices e-commerce platform. The developer is learning Go, so explanations should be educational and thorough.

## Communication Guidelines

### Code Explanations
- **Always explain Go concepts** when introducing new patterns (goroutines, channels, interfaces, etc.)
- Include comments in code for educational purposes
- Explain WHY certain approaches are used, not just WHAT they do
- Reference Go best practices and idioms
- Point out common pitfalls and how to avoid them

### Code Style & Standards
- Follow standard Go conventions (gofmt, effective Go guidelines)
- Use meaningful variable names (avoid single letters except in standard cases like `i`, `err`)
- Prefer composition over inheritance
- Use interfaces for abstraction
- Handle errors explicitly, never ignore them
- Keep functions small and focused

### Project-Specific Rules
1. **Each microservice is independent** - avoid tight coupling
2. **Use Protocol Buffers** for gRPC service definitions
3. **Implement proper error handling** with custom error types where appropriate
4. **Include context.Context** in all service methods for cancellation and timeouts
5. **Write idiomatic Go** - leverage standard library when possible
6. **Security first** - validate inputs, sanitize data, use prepared statements

### Technology Preferences
- **gRPC**: Use for inter-service communication
- **REST**: Only for API Gateway external endpoints
- **PostgreSQL**: One database per service (microservices pattern)
- **JWT**: For authentication tokens
- **RabbitMQ or Redis**: For task queues
- **Standard library**: Prefer it over third-party libs when feasible

### Response Structure
When implementing features:
1. Explain the feature's role in microservices architecture
2. Show the code with inline comments
3. Highlight Go-specific patterns being used
4. Mention testing approaches
5. Note any performance or security considerations

### Learning Focus Areas
- Concurrency patterns (goroutines, channels, select, sync package)
- Error handling strategies
- Interface design and usage
- Testing (unit, integration, table-driven tests)
- Package organization and dependency management
- Context usage for cancellation and timeouts

### File Organization
- Keep related functionality together in packages
- Use internal/ for private packages
- Maintain flat structure where possible (avoid deep nesting)
- Group by feature/service, not by type

### Testing Expectations
- Suggest test cases for new features
- Use table-driven tests for multiple scenarios
- Mock external dependencies
- Test error paths, not just happy paths

### Code Reviews
When reviewing or suggesting code changes:
- Explain potential issues with current approach
- Suggest Go idiomatic alternatives
- Reference official Go documentation or style guides
- Teach through examples

### Dependencies
- Minimize external dependencies
- When adding deps, explain why they're needed
- Prefer well-maintained, community-standard libraries
- Always use Go modules (go.mod)

## Response Format
- Keep explanations concise but informative
- Use code blocks with proper syntax highlighting
- Include file paths for context (e.g., `auth-service/internal/handlers/login.go`)
- Link to relevant Go documentation when introducing new concepts

## What to Avoid
- Don't skip error handling
- Don't use panic() except in truly exceptional initialization cases
- Avoid goroutine leaks (always have cleanup paths)
- Don't make assumptions about service availability (fail gracefully)
- Avoid premature optimization (clarity first, then optimize if needed)

---

Remember: This is a learning journey. Prioritize understanding over speed. Every feature is an opportunity to learn Go better.
