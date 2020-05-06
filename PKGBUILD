# Maintainer: Sighery
pkgname=email-notification
pkgver=0.1.0
pkgrel=1
pkgdesc="Send emails using Gmail's API with OAuth2"
url='https://github.com/Sighery/email-notification'
arch=('x86_64')
license=('MIT')
makedepends=('go')
source=(
	'main.go'
	'LICENSE'
)
sha256sums=(
	'948ebcb05c02d8bab7504967b3d93611340d4b250bdcc9f873f42489366673e5'
	'1072753cc74bf1991f606dbe5efdda157f30f30ddd4c51883752c93814fe3ba8'
)

build() {
	go build -o $pkgname
}

package() {
	install -Dm755 $srcdir/$pkgname $pkgdir/usr/bin/$pkgname

	# MIT license needs to be installed separately
	install -Dm644 $srcdir/LICENSE $pkgdir/usr/share/licenses/$pkgname/LICENSE
}
