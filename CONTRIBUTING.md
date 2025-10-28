# Contributing to the URP Educational Repository

Thank you for your interest in improving this educational resource! This guide will help you contribute effectively.

## Ways to Contribute

### 1. Improve Documentation
- Fix typos or clarify confusing explanations
- Add diagrams or illustrations
- Translate content to other languages
- Add more detailed comments to code

### 2. Add Examples
- Implement URP in additional programming languages
- Create visualization tools
- Build interactive demos
- Add step-by-step tutorials

### 3. Enhance the Code
- Optimize performance
- Add new features (e.g., congestion control)
- Improve error handling
- Add more comprehensive logging

### 4. Expand Testing
- Add more test cases
- Create benchmarking tools
- Build network simulation scenarios
- Add performance analysis

## Code Structure

```
UDPAssignmentFun/
├── cmd/                    # Executable entry points
│   ├── sender/            # Sender main program
│   └── receiver/          # Receiver main program
├── internal/              # Core implementation
│   ├── urp/              # Protocol definitions (segment, checksum)
│   ├── protocol/         # State machines (sender & receiver)
│   ├── plc/              # Packet Loss & Corruption simulator
│   └── logger/           # Event logging
├── chapters/             # Educational content (step-by-step)
├── examples/             # Code examples in various languages
├── data/                 # Test files and logs
└── tests/                # Test scripts and utilities
```

## Development Workflow

### Setting Up Your Environment

1. **Fork and clone:**
   ```bash
   git clone https://github.com/YOUR_USERNAME/reliable-udp-implementation.git
   cd reliable-udp-implementation
   ```

2. **Install Go 1.22+:**
   - Download from [go.dev/dl](https://go.dev/dl/)

3. **Build and test:**
   ```bash
   go build ./cmd/sender
   go build ./cmd/receiver
   ./run_tests.ps1 -Quick  # Windows
   # or
   ./run_tests.sh          # Linux/Mac
   ```

### Making Changes

1. **Create a feature branch:**
   ```bash
   git checkout -b feature/your-feature-name
   ```

2. **Make your changes:**
   - Follow existing code style
   - Add comments for complex logic
   - Update documentation if needed

3. **Test your changes:**
   ```bash
   # Build
   go build ./...
   
   # Run tests
   ./run_tests.ps1
   
   # Check for errors
   go vet ./...
   go fmt ./...
   ```

4. **Commit with clear messages:**
   ```bash
   git add .
   git commit -m "feat: add visualization for sliding window"
   # or
   git commit -m "docs: clarify checksum algorithm explanation"
   # or
   git commit -m "fix: correct sequence number wrapping"
   ```

### Commit Message Guidelines

Use conventional commits:
- `feat:` - New feature
- `fix:` - Bug fix
- `docs:` - Documentation changes
- `test:` - Adding or fixing tests
- `refactor:` - Code refactoring
- `style:` - Code style changes (formatting)
- `perf:` - Performance improvements
- `chore:` - Build process or tooling changes

### Pull Request Process

1. **Push your branch:**
   ```bash
   git push origin feature/your-feature-name
   ```

2. **Create a Pull Request on GitHub**
   - Describe what you changed and why
   - Reference any related issues
   - Include screenshots/examples if applicable

3. **Code Review:**
   - Address feedback from reviewers
   - Keep discussions respectful and constructive

4. **Merge:**
   - Once approved, a maintainer will merge your PR

## Coding Standards

### Go Code Style

- Follow [Effective Go](https://go.dev/doc/effective_go)
- Use `gofmt` for formatting
- Add meaningful comments for exported functions
- Keep functions focused and small

**Example:**
```go
// ComputeChecksum calculates an 8-bit checksum by summing all bytes
// in the header and payload, then taking modulo 256.
func ComputeChecksum(data []byte) uint8 {
    var sum uint32
    for _, b := range data {
        sum += uint32(b)
    }
    return uint8(sum % 256)
}
```

### Documentation Style

- Use clear, simple language
- Include examples where helpful
- Add diagrams for complex concepts
- Write for beginners (assume minimal networking knowledge)

### Test Additions

If adding new tests:
- Document what the test validates
- Use descriptive test names
- Include expected outcomes
- Handle edge cases

## Educational Content Guidelines

### Writing Chapter Content

1. **Start with the "Why"** - Explain the problem before the solution
2. **Use Analogies** - Make networking concepts relatable
3. **Include Visual Aids** - Diagrams, flowcharts, packet illustrations
4. **Provide Examples** - Show code snippets with explanations
5. **End with Exercises** - Let learners practice

### Example Structure

```markdown
# Chapter X: [Title]

## The Problem
[Explain what challenge we're solving]

## The Solution
[Describe the approach]

## How It Works
[Step-by-step breakdown with diagrams]

## Code Implementation
[Annotated code examples]

## Try It Yourself
[Exercises and experiments]
```

## Adding a New Programming Language

To add URP implementation in a new language:

1. Create directory: `examples/[language]/`
2. Implement basic sender and receiver
3. Add README with:
   - Setup instructions
   - Dependencies
   - Running examples
   - Language-specific notes
4. Keep it simple - match the educational level

## Questions?

- Open an [Issue](https://github.com/MegaPanchamZ/reliable-udp-implementation/issues) for questions
- Start a [Discussion](https://github.com/MegaPanchamZ/reliable-udp-implementation/discussions) for ideas
- Tag maintainers for urgent matters

## Code of Conduct

- Be respectful and inclusive
- Welcome beginners and encourage learning
- Provide constructive feedback
- Focus on education and collaboration

## Recognition

Contributors will be acknowledged in:
- README.md contributors section
- Release notes
- Documentation credits

Thank you for helping make network programming education better! 🚀
