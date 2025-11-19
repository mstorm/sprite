# sprite

SVG 파일들을 하나의 sprite sheet로 합쳐주는 도구입니다.

## Features

- 여러 SVG 파일을 하나의 sprite sheet로 결합
- 1x, 2x, 3x 해상도 PNG 파일 자동 생성
- JSON 맵 파일 생성 (각 sprite의 위치와 크기 정보)
- resvg 기반의 고품질 SVG → PNG 변환

## Prerequisites

### resvg 설치

이 프로젝트는 SVG를 PNG로 변환하기 위해 [resvg](https://github.com/RazrFalcon/resvg)를 사용합니다.

#### Option 1: Cargo로 설치 (권장)

```bash
# Rust가 설치되어 있지 않다면
curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh

# resvg 설치
cargo install resvg
```

#### Option 2: 사전 빌드된 바이너리 다운로드

[resvg releases](https://github.com/RazrFalcon/resvg/releases) 페이지에서 운영체제에 맞는 바이너리를 다운로드하고 PATH에 추가하세요.

#### Option 3: Makefile 사용

```bash
make install-deps
```

## Installation

```bash
# 의존성 설치 및 빌드
make all

# 또는 수동으로
go build -o sprite .
```

## Usage

```bash
./sprite <svg-file-1> <svg-file-2> ...
```

### Example

```bash
./sprite testdata/dot.svg testdata/restaurant.svg testdata/mountain.svg testdata/airport.svg
```

생성되는 파일:
- `sprites.svg` - 결합된 SVG 파일
- `sprites.png` - 1x PNG sprite sheet
- `sprites@2x.png` - 2x PNG sprite sheet
- `sprites@3x.png` - 3x PNG sprite sheet
- `sprites.json` - 1x sprite 맵
- `sprites@2x.json` - 2x sprite 맵
- `sprites@3x.json` - 3x sprite 맵

## Development

### Build

```bash
make build
```

### Test

```bash
make test
```

### Run with test data

```bash
make run
```

### Clean

```bash
make clean
```

## License

MIT
