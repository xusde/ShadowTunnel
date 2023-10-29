package old

import (
	"fmt"
	"io"
	"libs/encryption"
	"net"
)

func SafeWrite(conn net.Conn, buf []byte, connType string, encryptType string, key string) (n int, err error) {
	if connType == "tcp" {
		if encryptType == "none" {
			n, err = conn.Write(buf)
			if err != nil {
				return n, err
			}
		} else if encryptType == "aes-256-gcm" {
			//fmt.Println("[write] buf:", buf)
			encryptedBytes := []byte(encryption.Encrypt(key, buf))
			outBuf := make([]byte, 2+len(encryptedBytes))
			// Convert to BigEndian
			outBuf[0] = byte(len(encryptedBytes) >> 8)
			outBuf[1] = byte(len(encryptedBytes))
			copy(outBuf[2:], encryptedBytes)
			n, err = conn.Write(outBuf)
			if err != nil {
				return n, err
			}
			//encryptedMsg := Encrypt(key, buf)
			//n, err = conn.Write([]byte(encryptedMsg))
			//if err != nil {
			//	return n, err
			//}
		}
	}
	return n, nil
}

func SafeRead(conn net.Conn, buf []byte, connType string, encryptType string, key string) (n int, err error) {
	if connType == "tcp" {
		if encryptType == "none" {
			n, err = conn.Read(buf)
			if err != nil {
				return n, err
			}
		} else if encryptType == "aes-256-gcm" {
			buffer := make([]byte, 1024)
			n, err = conn.Read(buffer)
			if err != nil {
				return n, nil
			}
			if n > 0 {
				decryptedMsg := encryption.Decrypt(key, string(buffer[2:n]))
				copy(buf, decryptedMsg)
				//fmt.Println("[read] buf:", buf)
				n = len(decryptedMsg)
			}

			//
			//tempBuf := make([]byte, 1024)
			//n, err = conn.Read(tempBuf)
			//if err != nil {
			//	return n, err
			//}
			//// Decrypt to buf
			//if n > 0 {
			//	decryptedMsg := Decrypt(key, string(tempBuf[0:n]))
			//	copy(buf, decryptedMsg)
			//	//fmt.Println("[read] buf:", buf)
			//	n = len(decryptedMsg)
			//}
		}
	}
	return n, nil
}

func SafeCopy(src net.Conn, dst net.Conn, connType string, encryptType string, key string, operation string) (written int64, err error) {
	if connType == "tcp" {
		if encryptType == "none" {
			//Copy(src, dst)
			written, err = io.Copy(dst, src)
			if err != nil {
				return written, err
			}
		} else {
			size := 0x3FFF
			inBuf := make([]byte, size)
			cache := make([]byte, 0)
			for {
				//nr, er := src.Read(inbuf)
				firstRead := true
				nr, er := src.Read(inBuf)

				if nr > 0 {
					if operation == "encrypt" {
						encryptedBytes := []byte(encryption.Encrypt(key, inBuf[0:nr]))
						outBuf := make([]byte, 2+len(encryptedBytes))
						// Convert to BigEndian
						outBuf[0] = byte(len(encryptedBytes) >> 8)
						outBuf[1] = byte(len(encryptedBytes))
						copy(outBuf[2:], encryptedBytes)

						fmt.Println("[safeCopy] out:", outBuf[0:2+len(encryptedBytes)])
						nw, ew := dst.Write(outBuf[0 : 2+len(encryptedBytes)])
						if nw > 0 {
							written += int64(nw)
						}
						if ew != nil {
							err = ew
							break
						}
						//cache = append(cache, inbuf[0:nr]...)
						//fmt.Println("[copy] out cache:", len(cache))
						//if nr < size {
						//	fmt.Println("[copy] out:", len(cache))
						//	outbuf = []byte(Encrypt(key, cache))
						//	cache = []byte{}
						//} else {
						//	continue
						//}

						// Encrypt to buf
						//fmt.Println("[copy] out:", nr)
						//outbuf[] = []byte(Encrypt(key, inbuf[0:nr]))
						//fmt.Println("[copy] out c:", len(outbuf))
					} else if operation == "decrypt" {
						var totalBytes int
						var totalBytesRead int
						outBuf := make([]byte, 0)
						if firstRead {
							firstRead = false
							// Convert to BigEndian
							totalBytes = int(inBuf[0])<<8 + int(inBuf[1])
							//copy(cache, inBuf[2:nr+2])
							cache = append(cache, inBuf[2:nr]...)
							totalBytesRead += nr
							if nr+2 < totalBytes {
								continue
							} else {
								decryptedBytes := encryption.Decrypt(key, string(cache[0:totalBytes]))
								outBuf = append(outBuf, decryptedBytes...)
								nw, ew := dst.Write(outBuf[0:len(decryptedBytes)])
								if nw > 0 {
									written += int64(nw)
								}
								if ew != nil {
									err = ew
									break
								}
								cache = []byte{}
							}
						} else {
							cache = append(cache, inBuf[0:nr]...)
							totalBytesRead += nr
							if totalBytesRead < totalBytes {
								continue
							} else {
								decryptedBytes := encryption.Decrypt(key, string(cache[0:totalBytes]))
								copy(outBuf, decryptedBytes)
								nw, ew := dst.Write(outBuf[0:len(decryptedBytes)])
								if nw > 0 {
									written += int64(nw)
								}
								if ew != nil {
									err = ew
									break
								}
								cache = []byte{}
							}

						}

						//cache = append(cache, inbuf[0:nr]...)
						//fmt.Println("[copy] in cache:", len(cache))
						//if nr < size {
						//	// Decrypt to buf
						//	//fmt.Println("[copy] in:", nr)
						//	//fmt.Println("cache:", cache)
						//	outbuf = Decrypt(key, string(cache))
						//	fmt.Println("[copy] in:", len(outbuf))
						//	cache = []byte{}
						//} else {
						//	//fmt.Println("cache:", cache)
						//	continue
						//}

						//// Decrypt to buf
						//fmt.Println("[copy] in c:", nr)
						//outbuf = Decrypt(key, string(inbuf[0:nr]))
						//fmt.Println("[copy] in:", len(outbuf))
					}
					//nw, ew := dst.Write(outBuf)
					//if nw > 0 {
					//	written += int64(nw)
					//}
					//if ew != nil {
					//	err = ew
					//	break
					//}
					//if nr != nw {
					//	err = io.ErrShortWrite
					//	break
					//}
				}
				if er != nil {
					if er != io.EOF {
						err = er
					}
					break
				}
			}
			return written, err
		}
	}
	return written, nil
}
