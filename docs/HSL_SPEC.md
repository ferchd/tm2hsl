# HSL Bytecode Specification v1.0

## Overview

HSL (High-Speed Lexing) is a binary format designed for efficient lexical analysis of programming languages. It represents a compiled TextMate grammar as deterministic finite automata optimized for fast execution.

## File Structure

```
HSL File Header (64 bytes)
├── Magic: "HSL1" (4 bytes)
├── Version: 1 (2 bytes)
├── Header Size: 64 (2 bytes)
├── String Table Offset (4 bytes)
├── Regex Table Offset (4 bytes)
├── Scope Table Offset (4 bytes)
├── State Table Offset (4 bytes)
├── Rule Table Offset (4 bytes)
├── Total File Size (4 bytes)
├── Checksum (4 bytes)
└── Flags (4 bytes)

String Table
├── Count (4 bytes)
├── Offsets (4 bytes × count)
└── Data (null-terminated strings)

Regex Table
├── Count (4 bytes)
└── Entries (variable)

Scope Table
├── Count (4 bytes)
└── Entries (variable)

State Table
├── Count (4 bytes)
└── Entries (variable)

Rule Table
├── Count (4 bytes)
└── Entries (variable)
```

## Tables

### String Table
Contains all string literals used in the bytecode:
- Regex patterns
- Scope names
- Token names

### Regex Table
Compiled regular expressions with precomputed bytecode.

### Scope Table
Hierarchical scope definitions for token classification.

### State Table
State machine states with transitions.

### Rule Table
Matching rules combining regexes, actions, and state transitions.

## Execution Model

1. Start in initial state
2. Match input against current state's rules
3. Apply actions (push/pop scopes, emit tokens)
4. Transition to next state
5. Repeat until end of input

## Compatibility

- Bytecode version 1 is backward compatible
- New features add optional tables
- Engines can ignore unknown sections