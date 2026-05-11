# Interview Prep Agent Instructions

## Project Purpose

This workspace is for job interview preparation. The user will keep adding interview questions, interview notes, job descriptions, and resume material. The agent should turn that raw material into organized, reviewable Markdown documents.

## Working Language

- Default to Chinese for explanations, answers, summaries, and review notes.
- Preserve original English terms from job descriptions, resumes, and technical keywords when they are important for interviews.
- Use concise, interview-ready wording. Avoid inflated claims, vague praise, or unsupported facts.
- Assume the user is a beginner for many technical topics. Explain prerequisites before advanced details.
- Prefer plain-language explanations first, then technical terms, then interview-ready wording.

## Core Responsibilities

When the user provides new material:

1. Classify it as one or more of:
   - interview question
   - interview transcript or notes
   - job description
   - resume content
   - project experience
   - knowledge topic
2. Extract reusable interview points:
   - key question
   - expected interviewer intent
   - standard answer
   - resume or project evidence
   - follow-up questions
   - weak spots to review
3. Produce or update reviewable Markdown files in the existing folder structure.
4. Keep the content practical for real interviews, not textbook-only.

## Answer Format For Interview Questions

Prefer this structure unless the user asks for another format:

```markdown
## 问题

...

## 面试官考察点

...

## 标准回答

...

## 结合我的经历

...

## 可展开追问

...

## 复习要点

...
```

For technical questions, include:

- prerequisite concepts
- plain-language explanation
- short answer first
- deeper explanation
- production scenario
- common pitfalls
- how to connect it to the user's projects or resume

For behavioral questions, include:

- STAR structure
- concise spoken version
- stronger resume-aligned version
- possible interviewer follow-ups

## Document Organization

Use the existing directory style:

- `interview/` for interview knowledge, questions, and topic notes.
- `resume/` for resume drafts, resume bullet polishing, and experience alignment.
- `projects/` for project architecture, project stories, and business background.
- `templates/` for reusable Markdown structures.
- `reviews/` for mock interview feedback and interview retrospectives.

When creating new files:

- Use lowercase, hyphenated filenames.
- Put topic notes under a relevant subdirectory.
- Add an index or README only when it helps navigation.
- Do not create unnecessary directories for one-off notes.

Use templates when applicable:

- `templates/interview-question.md` for a single interview question.
- `templates/mysql-question.md` for MySQL interview questions and MySQL 八股整理.
- `templates/interview-note.md` for interview notes or retrospectives.
- `templates/jd-resume-match.md` for JD and resume alignment.
- `templates/project-story-star.md` for project experience stories.

## Quality Rules

- Do not fabricate project facts, metrics, company names, or production incidents.
- If a detail is missing, mark it clearly as `待补充`.
- When improving resume or project wording, keep it truthful and defensible in follow-up questions.
- Prefer answers that a candidate can speak in 1-2 minutes, then add optional deeper points.
- Highlight risky or weak claims that may be challenged by an interviewer.
- Keep Markdown clean and easy to review.
- Do not assume the user already understands jargon. Define terms like `索引`, `事务`, `锁`, `Buffer Pool`, `MVCC`, `redo log`, and `binlog` when first used in a note.
- For hard concepts, use a three-layer structure: `先讲人话` -> `再讲原理` -> `面试怎么说`.

## When Updating Existing Notes

- Read nearby existing files first and follow their style.
- Preserve useful existing content.
- Merge duplicate questions instead of scattering repeated answers.
- If the same topic appears in multiple places, prefer improving the most relevant existing file.

## User Collaboration

- If the user sends raw notes, organize them without requiring perfect formatting.
- If the user asks for standard answers, provide polished answers directly.
- If the user asks for interview simulation, act as the interviewer first, then give feedback after the answer.
- If the user asks for resume alignment with a JD, map JD requirements to resume evidence and identify gaps.
