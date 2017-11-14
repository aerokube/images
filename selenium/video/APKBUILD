# Contributor: Anton Dzyk <anton.dzyk@yandex.ru>
# Maintainer: Anton Dzyk <anton.dzyk@yandex.ru>
pkgname=ffmpeg-x11grab
pkgsource=ffmpeg
pkgver=3.4
pkgrel=1
pkgdesc="Save stream video on xvfb to h264 video"
url="http://ffmpeg.org/"
arch="all"
license="GPL"
makedepends="gnutls-dev lame-dev libvorbis-dev xvidcore-dev zlib-dev libvdpau-dev
	imlib2-dev x264-dev libtheora-dev coreutils bzip2-dev perl-dev libvpx-dev
	libvpx-dev sdl2-dev libxfixes-dev libva-dev alsa-lib-dev rtmpdump-dev
	v4l-utils-dev yasm opus-dev x265-dev"
source="http://ffmpeg.org/releases/ffmpeg-$pkgver.tar.xz
	0001-libavutil-clean-up-unused-FF_SYMVER-macro.patch
	"
builddir="$srcdir/$pkgsource-$pkgver"

build() {
	local _dbg="--disable-debug"
	local _asm=""
	[ -n "$DEBUG" ] && _dbg="--enable-debug"

	case "$CARCH" in
	x86 | arm*) _asm="--disable-asm" ;;
	esac

	cd "$builddir"
	./configure \
		--prefix=/usr \
		--enable-avfilter \
		--enable-gpl \
		--enable-libx264 \
		--disable-stripping \
		--disable-ffplay \
		--disable-ffprobe \
		--disable-ffserver \
		--disable-doc \
		--disable-htmlpages \
		--disable-manpages \
		--disable-podpages \
		--disable-txtpages \
		--disable-w32threads \
		--disable-alsa \
		--disable-audiotoolbox \
		--disable-cuda \
		--disable-cuvid \
		--disable-d3d11va \
		--disable-dxva2 \
		--disable-nvenc \
		--disable-vaapi \
		--disable-vda \
		--disable-vdpau \
		--disable-videotoolbox \
		$_asm $_dbg
	make
}

package() {
	cd "$builddir"
	make DESTDIR="$pkgdir" install
}

libs() {
	pkgdesc="Libraries for ffmpeg"
	replaces="ffmpeg"
	mkdir -p "$subpkgdir"/usr
	mv "$pkgdir"/usr/lib "$subpkgdir"/usr
}

sha512sums="357445f0152848d43f8a22f1078825bc44adacff9194e12cc78e8b5edac8e826bbdf73dc8b37e0f2a3036125b76b6b9190153760c761e63ebd2452a39e39536f  ffmpeg-3.4.tar.xz
32652e18d4eb231a2e32ad1cacffdf33264aac9d459e0e2e6dd91484fced4e1ca5a62886057b1f0b4b1589c014bbe793d17c78adbaffec195f9a75733b5b18cb  0001-libavutil-clean-up-unused-FF_SYMVER-macro.patch"
