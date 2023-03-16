package generator

import (
	"fmt"
	"math"
	"strings"
)

func (i *Insn) genCodeVdRs1mVs2Vm() []string {
	nfields := getNfieldsRoundedUp(i.Name)
	combinations := i.combinations(
		nfieldsLMULs(nfields),
		[]SEW{getEEW(i.Name)},
		[]bool{false, true})
	res := make([]string, 0, len(combinations))

	for _, c := range combinations {
		builder := strings.Builder{}
		builder.WriteString(c.comment())

		builder.WriteString(i.gWriteRandomData(LMUL(1)))
		builder.WriteString(i.gLoadDataIntoRegisterGroup(0, LMUL(1), SEW(32)))

		lmul1 := LMUL(math.Max(float64(c.LMUL)*float64(nfields), 1))
		vd := int(lmul1)
		vs2 := 2 * int(c.LMUL1) * nfields
		builder.WriteString(i.gWriteIndexData(lmul1, c.Vl, c.SEW))
		builder.WriteString(i.gLoadDataIntoRegisterGroup(vs2, lmul1, c.SEW))
		builder.WriteString(i.gWriteRandomData(lmul1))
		builder.WriteString(i.gLoadDataIntoRegisterGroup(vd, lmul1, c.SEW))
		builder.WriteString(i.gWriteIntegerTestData(lmul1, c.SEW, 0))

		builder.WriteString("# -------------- TEST BEGIN --------------\n")
		builder.WriteString(i.gVsetvli(c.Vl, c.SEW, c.LMUL))
		builder.WriteString(fmt.Sprintf("%s v%d, (a0), v%d%s\n", i.Name, vd, vs2, v0t(c.Mask)))
		builder.WriteString("# -------------- TEST END   --------------\n")

		builder.WriteString(i.gResultDataAddr())
		builder.WriteString(i.gStoreRegisterGroupIntoResultData(vd, lmul1, c.SEW))
		builder.WriteString(i.gMagicInsn(vd))

		res = append(res, builder.String())
	}

	return res
}