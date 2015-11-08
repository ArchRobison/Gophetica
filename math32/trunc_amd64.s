#include "textflag.h"

// func Trunc(x float32) float32
TEXT ·Trunc(SB), NOSPLIT, $0
    MOVQ x+0(FP), X0
    BYTE $0xc4; BYTE $0xe3; BYTE $0x79; BYTE $0x0a; BYTE $0xc0; BYTE $0x03 // ROUNDSS X0, X0, 3
    MOVQ X0, ret+8(FP)
    RET
