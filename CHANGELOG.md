# Changelog

All notable changes to MaestroSDK will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

---

## [Unreleased]

### Fixed

- Fix `ApplyList` error reporting — now includes resource kind, name, and failure reason for each failed item instead of aggregating by count only; returns structured per-item error output (`resource/list.go`) (#152)
