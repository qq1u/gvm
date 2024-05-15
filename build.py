import os
import shutil
import argparse
import subprocess
from pathlib import Path
from datetime import datetime

current_dir = Path(__file__).parent
dist_dir = current_dir / "dist"


class System:
    WINDOWS = "windows"
    DARWIN = "darwin"
    LINUX = "linux"


class Arch:
    AMD64 = "amd64"
    ARM64 = "arm64"
    X86 = "386"


def mkdir(path):
    os.makedirs(path, exist_ok=True)


def build(tag, system, arch):
    start = datetime.now()
    print(f"Start build {tag} {system} {arch} at {start}")
    target_dir = dist_dir / f"gvm_{system}_{arch}" / "gvm"  # gvm0.0.1_windows_amd64/gvm/gvm.exe
    mkdir(target_dir)
    env = os.environ.copy()
    env.update({"GOOS": system, "GOARCH": arch, "CGO_ENABLED": "0"})
    output = target_dir / "gvm.exe" if system == System.WINDOWS else target_dir / "gvm"

    cwd = current_dir / "gvm" / "cmd"
    cmd = [
        "go",
        "build",
        f'-ldflags=-w -s -X gvm.VERSION={tag}',
        "-o",
        str(output),
    ]
    subprocess.run(cmd, env=env, cwd=cwd)

    shutil.make_archive(str(target_dir.parent), "zip", target_dir.parent)
    shutil.rmtree(target_dir.parent)
    print(f"Finish build {tag} {system} {arch} used: {(datetime.now() - start).total_seconds()}s\n")


def main():
    result = subprocess.run(["git", "describe", "--tags", "--abbrev=0"], capture_output=True, text=True)
    tag = result.stdout.strip()

    parser = argparse.ArgumentParser()
    parser.add_argument("--os", type=str, default=None)
    parser.add_argument("--arch", type=str, default=None)
    kwargs = parser.parse_args()

    shutil.rmtree(dist_dir, ignore_errors=True)

    system_list = [System.WINDOWS, System.DARWIN, System.LINUX]
    os_ = kwargs.os
    if os_ and os_ in system_list:
        system_list = [os_]

    arch_list = [Arch.AMD64, Arch.X86, Arch.ARM64]
    arch = kwargs.arch
    if arch and arch in arch_list:
        arch_list = [arch]

    for system in system_list:
        for arch in arch_list:
            if system == System.DARWIN and arch == Arch.X86:
                continue
            build(tag, system, arch)


if __name__ == "__main__":
    main()
