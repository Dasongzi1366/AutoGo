
# AutoGo —— 高性能 Android 自动化测试框架

**AutoGo** 是一个基于 Go 语言开发的 Android 自动化操作工具，旨在为用户提供一个高效、安全、灵活的 Android 自动化解决方案。与传统的自动化工具不同，AutoGo 的二进制执行文件能够在 Android 系统上直接运行，具备强大的跨应用操作能力，同时无需安装任何 APK。这使得它在复杂的自动化场景中，尤其是在需要与其他应用无缝交互时，具有独特的优势。

## 功能亮点

- **非侵入式架构**：AutoGo 不依赖于任何 APK 安装，通过 ADB 或 root 权限的 shell 即可在 Android 系统中运行，无需修改系统或应用。
  
- **强大的自动化能力**：AutoGo 支持对 Android 系统和应用的全面控制，涵盖从 UI 操作到系统层级的复杂交互。支持以下模块：
  - **accessibility**：无障碍辅助控制功能。
  - **app**：管理与操作 Android 应用的生命周期。
  - **device**：获取设备信息及执行系统操作。
  - **files**：文件系统管理，便捷地读取、写入和修改文件。
  - **frpc**：集成网络穿透服务，提供远程访问能力。
  - **https**：安全的 HTTP/HTTPS 请求支持。
  - **ime**：输入法和文本输入的自动化控制。
  - **imgui**：图形用户界面交互支持。
  - **keys**：模拟按键和输入操作。
  - **media**：多媒体操作管理。
  - **memory**：内存管理和数据读写。
  - **ppocr**：集成光学字符识别 (OCR) 功能，轻松识别屏幕文本。
  - **storages**：存储操作支持。
  - **system**：系统级别的命令执行与操作。
  - **touch**：精准的触摸屏幕控制。
  - **turocr**：高级 OCR 识别与解析。
  - **yolo**：目标检测与计算机视觉处理。

- **可扩展性强**：开发者可以通过 Go 语言编写自定义脚本，利用丰富的内置 API 实现高度定制化的自动化流程。

- **高安全性**：AutoGo 编译为 Android 可执行的二进制文件，极大地提高了程序的安全性和防破解能力。相比传统的脚本工具，AutoGo 不暴露源码，有效防止逆向工程。

- **集成与兼容性**：支持在任意 Android 应用中通过 root 权限集成运行，无论是作为独立工具使用，还是作为其他应用的一部分嵌入，均能发挥其强大的自动化功能。

## 为什么选择 AutoGo？

- **跨平台支持**：AutoGo 提供的二进制文件可以在不同 Android 设备上运行，无需 APK 安装，使用 ADB 即可完成自动化操作。
  
- **灵活的开发模式**：用户通过 Go 语言编写自动化脚本，充分利用 Go 语言的高并发和高效特性，快速实现复杂的自动化任务。

- **安全与隐私**：通过编译为二进制文件，AutoGo 能够有效防止源码泄露和逆向工程，确保自动化脚本的安全性。

- **丰富的功能模块**：从应用层的辅助功能控制到系统级别的命令执行，AutoGo 涵盖了 Android 系统操作的方方面面，几乎可以满足所有自动化场景的需求。

## 适用场景

- **应用自动化测试**：开发者可以通过 AutoGo 编写自动化脚本，测试应用的功能、性能和兼容性。
- **跨应用操作**：实现应用之间的无缝交互操作，比如自动化测试跨多个应用的场景。
- **高安全性要求的自动化**：需要防止被逆向破解的自动化脚本，例如涉及敏感操作的自动化流程。
- **内置集成**：AutoGo 可以嵌入到其他 Android 应用中，提供自动化能力，为复杂的应用开发提供强有力的支持。

## 安装指南

1. 运行以下命令安装 AutoGo：
   ```bash
   go install github.com/Dasongzi1366/AutoGoBuild@latest && AutoGoBuild
   ```

2. 下载并配置 `ANDROID_NDK_HOME` 环境变量：

   你可以从以下链接下载 NDK r21e：
   [NDK r21e 下载地址 (Windows)](https://dl.google.com/android/repository/android-ndk-r21e-windows-x86_64.zip)

   下载完成后，解压到你希望的路径，并将 `ANDROID_NDK_HOME` 设置为该路径：
   ```bash
   set ANDROID_NDK_HOME=C:\android-ndk
   ```

3. 编译二进制文件：

   **编译 arm64 架构：**
   ```bash
   AutoGo build arm64
   ```

   **编译 x86_64(模拟器) 架构：**
   ```bash
   AutoGo build x86_64
   ```

4. 运行流程：
   - 推送二进制文件到 Android 设备：
     ```bash
     adb push [name] /data/local/tmp/app
     ```
   - 设置执行权限：
     ```bash
     adb shell chmod 755 /data/local/tmp/app
     ```
   - 执行二进制文件：
     ```bash
     adb shell /data/local/tmp/app
     ```

## 示例代码

以下是一个简单的 AutoGo 示例，展示了如何使用无障碍辅助功能查找文本元素：

```go
package main

import (
    "fmt"
    "github.com/Dasongzi1366/AutoGo/accessibility"
)

func main() {
    a := accessibility.New()
    obj := a.Text("应用商店").FindOnce()
    if obj != nil {
        fmt.Println(obj.ToString())
    } else {
        fmt.Println("没有找到")
    }
}
```

## 联系方式

如需更多信息或获取技术支持，欢迎通过以下方式与我们联系：
- **QQ群**: 753399754
