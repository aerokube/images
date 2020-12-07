# Contributor: Anton Dzyk <anton.dzyk@yandex.ru>
# Maintainer: Ivan Krutov <vania-pooh@aerokube.com>
pkgname=ffmpeg-x11grab
pkgsource=ffmpeg
pkgver=4.3
pkgrel=1
pkgdesc="FFMPEG version optimized for x11grab to x264"
url="https://ffmpeg.org"
arch="all"
license="LGPL-2.1-or-later GPL-2.0-or-later"
makedepends="gnutls-dev lame-dev libvorbis-dev xvidcore-dev zlib-dev libvdpau-dev
	imlib2-dev x264-dev libtheora-dev coreutils bzip2-dev perl-dev libvpx-dev
	libvpx-dev sdl2-dev libxfixes-dev libva-dev v4l-utils-dev yasm opus-dev libass-dev pulseaudio-dev"
source="https://ffmpeg.org/releases/ffmpeg-$pkgver.tar.xz
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
		--enable-libpulse \
		--enable-static \
		--enable-small \
		--disable-ffplay \
		--disable-ffprobe \
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
		--disable-vdpau \
		--disable-videotoolbox \
		--disable-librtmp \
		--disable-devices \
		--enable-indev=xcbgrab \
		--enable-indev=pulse \
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

sha512sums="f031eb6c4423887af323ab7d1f431234d4e30993a52db45dccf427b41eb442a3bd020dcbc13e83cbf813fad0f36c849cb651203570148387c864507aa19f313a  ffmpeg-4.3.tar.xz
1047a23eda51b576ac200d5106a1cd318d1d5291643b3a69e025c0a7b6f3dbc9f6eb0e1e6faa231b7e38c8dd4e49a54f7431f87a93664da35825cc2e9e8aedf4  0001-libavutil-clean-up-unused-FF_SYMVER-macro.patch"
