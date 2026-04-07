# Changelog

All notable changes to MaestroSDK will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

---

## [Unreleased]

---

## [0.1.4] - 2026-04-06

### Fixed

- fix: swap Credential/Workspace in DependencyOrder — credentials now apply after workspaces, fixing restore of workspace-scoped credentials (`resource/list.go`) (#195)
