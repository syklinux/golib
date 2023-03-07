package collection

import (
	"os"
	"syscall"
)

func FileLockWrite(localPath string, content []byte) (err error) {
	f, err := os.OpenFile(localPath, os.O_WRONLY|syscall.O_CREAT|os.O_TRUNC, 0644)
	if err != nil {
		return err
	} else {
		err = syscall.Flock(int(f.Fd()), syscall.LOCK_EX|syscall.LOCK_NB)
		if err != nil {
			return err
		} else {
			if _, err = f.Write(content); err != nil {
				return err
			}
			syscall.Flock(int(f.Fd()), syscall.LOCK_UN)
		}
		f.Close()
	}
	return err
}
