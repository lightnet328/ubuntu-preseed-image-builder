package main

import (
	"fmt"
	"os"
	"os/exec"
)

func exist(path string) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		panic(err)
	}
}

func command(c string) {
	if err := exec.Command("sh", "-c", c).Run(); err != nil {
		panic(err)
	}
}

func write(path string, s string) {
	f, err := os.OpenFile(path, os.O_WRONLY, 0600)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	fmt.Fprintln(f, s)
}

const isolinux = `
default install
label install
  menu label ^Install Ubuntu Server
  kernel /install/vmlinuz
  append DEBCONF_DEBUG=5 auto=true locale=en_US.UTF-8 console-setup/charmap=UTF-8 console-setup/layoutcode=us console-setup/ask_detect=false pkgsel/language-pack-patterns=pkgsel/install-language-support=false file=/cdrom/preseed/custom.seed vga=normal initrd=/install/initrd.gz quiet --
label hd
  menu label ^Boot from first hard disk
  localboot 0x80
`

const grub = `
menuentry "Install Ubuntu Server" {
  set gfxpayload=keep
  linux  /install/vmlinuz file=/cdrom/preseed/custom.seed debian-installer/locale=en_US console-setup/layoutcode=us quiet ---
  initrd /install/initrd.gz
}
`

func regenerateISO(image, suffix string) {
	if suffix != "" {
		suffix = "-" + suffix
	}
	input := fmt.Sprintf("%s.iso", image)
	output := fmt.Sprintf("%s%s.iso", image, suffix)
	exist(fmt.Sprintf("/builder/%s", input))
	command(fmt.Sprintf("mount -t iso9660 /builder/%s /media", input))
	command("cd /media && find . ! -type l | cpio -pdum /ubuntu")
	command("cp /builder/preseed.cfg /ubuntu/preseed/custom.seed")
	write("/ubuntu/isolinux/isolinux.cfg", isolinux)
	write("/ubuntu/boot/grub/grub.cfg", grub)
	command(fmt.Sprintf("cd /ubuntu && xorriso -as mkisofs -l -J -R -b isolinux/isolinux.bin -c isolinux/boot.cat -no-emul-boot -boot-load-size 4 -boot-info-table -o /builder/%s /ubuntu", output))
	command(fmt.Sprintf("isohybrid /builder/%s", output))
}
