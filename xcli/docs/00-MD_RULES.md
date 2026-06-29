MDRULES.md
================================

## 1. PURPOSE & PRINCIPLE

This document establishes the structural formatting rules for all Markdown (.md) documentation within this repository (and only MarkDown!). 
The absolute driving principle is "Human Readability First". Documents must be designed to be comfortably scanned and read by human eyes directly in their raw text or code editor format, without relying on HTML rendering engines.


&nbsp;
&nbsp;


________________________________________________________________________________

## 2. HORIZONTAL SEPARATORS & MAJOR RESPIROS (H2)

Major thematic transitions must be separated using explicit visual anchors to prevent blocky text accumulation.

* Every major Section Header (H2) must be preceded by exactly:
  - Two empty lines
  - One `&nbsp;` line
  - One `&nbsp;` line
  - Two empty lines
  - A full horizontal rule line containing 80 underscores (`________________________________________________________________________________`)
  - One empty line

* The actual text content or paragraphs directly below an H2 header must always be preceded by exactly one empty line.


&nbsp;


### Visual Example for H2 Structure:

```text


&nbsp;
&nbsp;


________________________________________________________________________________

## TARGET H2 HEADER

This is the content text starting after exactly one empty line.
```


&nbsp;
&nbsp;


________________________________________________________________________________

## 3. VERTICAL SPACING BY HEADER HIERARCHY

To ensure proper visual breathing space, strict empty line structures must be applied before and after headers depending on their level.


&nbsp;


### 3.1 Sub-section Headers (H3)

* Preceding space: Must have exactly two empty lines, followed by one line with `&nbsp;`, followed by exactly two empty lines before the H3 header.
* Proceeding space: Must have exactly one empty line before the content begins.

&nbsp;

#### Visual Example for H3 Structure:

```text


&nbsp;


### TARGET H3 HEADER

This is the sub-section content starting after exactly one empty line.
```


&nbsp;


### 3.2 Deep Headers (H4)

* Preceding space: Must have exactly one empty line, followed by one line with `&nbsp;`, followed by exactly one empty line before the H4 header.
* Proceeding space: Must have exactly one empty line of spacing before the text or content begins.

&nbsp;

#### Visual Example for H4 Structure:
```text

&nbsp;

#### TARGET H4 HEADER

This is the deep header content starting after exactly one empty line.
```


&nbsp;


### 3.3 Minor Headers (H5 and H6)

* Preceding space: Must have at least one empty line before the header.
* Proceeding space: Post-header spacing is optional (facultativo) and left to the author's discretion.


&nbsp;
&nbsp;


________________________________________________________________________________

## 4. LISTS AND BULLET POINTS

* List items must never be cramped together.
* Bullet items must have a clean indent and use short, punchy fragments instead of dense paragraphs.
* Sub-bullets or sub-lists must follow the same breathing pattern without cluttering the vertical alignment.


&nbsp;
&nbsp;


________________________________________________________________________________

## 5. CODE BLOCKS EMBEDDING

* All code blocks must be cleanly isolated.
* Code blocks must have exactly one empty line before the starting triple backticks (```) and one empty line after the closing triple backticks.
* File types or language names should be explicitly stated next to the opening backticks when relevant (e.g., ```example.go) to preserve context in raw view.
