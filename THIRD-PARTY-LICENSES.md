# Third-Party Licenses

This document lists the open-source dependencies used by Stardoc and their respective licenses.

## Stardoc License

Stardoc itself is licensed under the [Beerware License](LICENSE). However, it depends on the following third-party software, each with its own license terms.

---

## Go Dependencies

### Direct Dependencies

#### 1. Cobra
- **Package**: `github.com/spf13/cobra` v1.10.2
- **License**: Apache License 2.0
- **Copyright**: © 2013 Steve Francia
- **License URL**: https://github.com/spf13/cobra/blob/main/LICENSE.txt
- **Purpose**: CLI framework for building command-line applications

#### 2. Godog
- **Package**: `github.com/cucumber/godog` v0.15.1
- **License**: MIT License
- **Copyright**: © SmartBear Software
- **License URL**: https://github.com/cucumber/godog/blob/main/LICENSE
- **Purpose**: BDD testing framework (dev dependency)

### Key Indirect Dependencies

#### 3. Lipgloss
- **Package**: `github.com/charmbracelet/lipgloss` v1.1.0
- **License**: MIT License
- **Copyright**: © 2021-2025 Charmbracelet, Inc.
- **License URL**: https://github.com/charmbracelet/lipgloss/blob/master/LICENSE
- **Purpose**: Terminal styling library

#### 4. Spinner
- **Package**: `github.com/briandowns/spinner` v1.23.2
- **License**: Apache License 2.0
- **Copyright**: © Brian J. Downs
- **License URL**: https://github.com/briandowns/spinner
- **Purpose**: Terminal progress indicators

#### 5. UUID
- **Package**: `github.com/google/uuid` v1.6.0
- **License**: BSD 3-Clause License
- **Copyright**: © 2009,2014 Google Inc.
- **License URL**: https://github.com/google/uuid/blob/master/LICENSE
- **Purpose**: Universally unique identifier generation

#### 6. Color
- **Package**: `github.com/fatih/color` v1.7.0
- **License**: MIT License
- **Copyright**: © 2013 Fatih Arslan
- **License URL**: https://github.com/fatih/color/blob/main/LICENSE.md
- **Purpose**: Terminal color output

#### 7. YAML
- **Package**: `gopkg.in/yaml.v3` v3.0.1
- **License**: Apache License 2.0 and MIT License
- **Copyright**: © 2011-2019 Canonical Ltd.
- **License URL**: https://github.com/go-yaml/yaml/blob/v3/LICENSE
- **Purpose**: YAML parsing and generation

---

## Node.js Dependencies

The following npm packages are installed as part of the Starlight documentation site that Stardoc generates:

### Runtime Dependencies

#### 1. Astro
- **Package**: `astro` ^5.0.0
- **License**: MIT License
- **Copyright**: © [Astro contributors]
- **License URL**: https://github.com/withastro/astro/blob/main/LICENSE
- **Purpose**: Static site generator framework
- **Website**: https://astro.build

#### 2. Starlight
- **Package**: `@astrojs/starlight` ^0.37.3
- **License**: MIT License
- **Copyright**: © [Astro contributors]
- **License URL**: https://github.com/withastro/starlight
- **Purpose**: Documentation theme for Astro
- **Website**: https://starlight.astro.build

#### 3. Sharp
- **Package**: `sharp` ^0.34.0
- **License**: Apache License 2.0
- **Copyright**: © 2013 Lovell Fuller and contributors
- **License URL**: https://github.com/lovell/sharp/blob/main/LICENSE
- **Purpose**: High-performance image processing
- **Website**: https://sharp.pixelplumbing.com

### Development Dependencies

#### 4. TypeScript
- **Package**: `typescript` ^5.7.3
- **License**: Apache License 2.0
- **Copyright**: © Microsoft Corporation
- **License URL**: https://github.com/microsoft/TypeScript/blob/main/LICENSE.txt
- **Purpose**: TypeScript language support
- **Website**: https://www.typescriptlang.org

#### 5. Node.js Type Definitions
- **Package**: `@types/node` ^22.10.5
- **License**: MIT License
- **Copyright**: © Microsoft Corporation
- **License URL**: https://github.com/DefinitelyTyped/DefinitelyTyped
- **Purpose**: TypeScript type definitions for Node.js

---

## License Compatibility

All third-party dependencies use permissive open-source licenses (MIT, Apache 2.0, BSD 3-Clause) that are compatible with the Beerware License used by Stardoc.

### License Summary

- **MIT License**: Very permissive; allows commercial use, modification, distribution, and private use. Requires only that the license and copyright notice be included.

- **Apache License 2.0**: Similar to MIT but includes explicit patent grant and requires stating significant changes to the code.

- **BSD 3-Clause**: Permissive license similar to MIT with an additional clause preventing use of contributors' names for endorsement.

- **Beerware License**: Extremely permissive; allows anything as long as you retain the notice. Optional: buy the author a beer if you meet and think it's worth it! 🍺

---

## Full License Texts

For the complete license texts of each dependency, please refer to the License URL provided for each package above.

## Compliance

By using Stardoc, you agree to comply with the license terms of all included third-party software. The copyright and license notices of these dependencies are preserved as required by their respective licenses.

---

**Note**: This list includes only the primary dependencies. Each of these packages may have their own dependencies with additional licenses. For a complete dependency tree and all transitive dependencies:

- **Go dependencies**: Run `go list -m all` in the project directory
- **Node.js dependencies**: Check `node_modules` after running `npm install` in a generated site

For any questions about licensing, please open an issue at: https://github.com/nicovandenhove/stardoc/issues

---

*Last updated: 2026-02-15*
