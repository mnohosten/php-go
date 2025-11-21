# RFCs (Request for Comments)

This directory contains design proposals for significant features or changes to PHP-Go.

## Purpose

RFCs are used for:

1. **Major Features** - New language features or extensions
2. **Breaking Changes** - Changes that affect compatibility
3. **Architecture Changes** - Significant internal changes
4. **New APIs** - Public API additions
5. **Complex Features** - Features needing community input

## RFC Process

### 1. Draft
- Create RFC document
- Explain motivation and use cases
- Propose solution with examples
- Discuss alternatives
- Outline implementation plan

### 2. Discussion
- Share with team/community
- Gather feedback
- Revise based on input
- Address concerns

### 3. Decision
- Accept: Proceed with implementation
- Reject: Archive with rationale
- Defer: Revisit later

### 4. Implementation
- Implement as specified
- Update documentation
- Add tests
- Mark RFC as completed

## RFC Template

```markdown
# RFC-NNNN: Title

## Status
Draft | Discussion | Accepted | Rejected | Implemented

## Summary
Brief description (2-3 sentences)

## Motivation
Why is this needed?

## Proposal
What exactly are we proposing?

## Examples
Code examples showing usage

## Implementation
How will this be implemented?

## Alternatives
What other approaches were considered?

## Impact
- Compatibility impact
- Performance impact
- Complexity impact

## Timeline
Proposed implementation schedule
```

## Active RFCs

### Accepted
- (None yet - project just starting)

### Under Discussion
- (None yet)

## Completed RFCs

- (None yet)

## Rejected RFCs

- (None yet)

## Topics for Future RFCs

### Language Features
- Custom parallel constructs syntax
- PHP-Go specific optimizations
- Go-style defer statements
- Pattern matching enhancements

### Integration Features
- Go struct ↔ PHP class mapping
- Go channel integration
- Go context integration
- Async/await syntax

### Performance Features
- JIT compiler design
- Advanced optimizer
- Profile-guided optimization
- Memory optimizations

### Ecosystem Features
- Package manager integration
- Extension marketplace
- IDE integration
- Debugging protocol

## Contributing RFCs

Anyone can propose an RFC:

1. Copy the RFC template
2. Fill out all sections thoughtfully
3. Create a pull request or issue
4. Engage in discussion
5. Iterate based on feedback

## Guidelines

**Good RFCs**:
- Clear motivation
- Concrete examples
- Thorough analysis
- Consider alternatives
- Implementation plan

**Avoid**:
- Vague proposals
- Missing use cases
- No alternatives discussed
- Unrealistic scope
- Breaking changes without strong justification

---

**Last Updated**: 2025-11-21
