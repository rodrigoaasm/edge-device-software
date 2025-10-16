package plc_runtime

import (
	"math"
	"math/rand/v2"
	"time"

	"github.com/go-logr/logr"
)

var (
	OPC_PATH_EXTERNAL        = "|var|CODESYS Control Win V3 x64.Application"
	OPC_NS_EXTERNAL   uint16 = 4
	OPC_NS_INTERNAL   uint16 = 1
)

type PLCRunner struct {
	stopCh              chan bool
	tWindow             int
	screwPositionTs     []float64
	injectScrewPosition int
	manu1ScrewPosition  int
	manu2ScrewPosition  int
	returnScrewPosition int
	machineInjectionTs  []float64
	startInjection      int
	maxStep1Injection   int
	maxStep2Injection   int
	repressInjection    int
	down1Pressure       int
	down2Pressure       int
	machinePressuseTs   []float64
	startMaxPressure    int
	endMaxPressure      int
	machineTempTs       []float64
	maxTemp             int
}

func NewPLCRunner(hertz int) *PLCRunner {
	v := 1 / float64(hertz)
	tWindow := int(v * 1000)

	rel := (float64(hertz*5) / 3000)

	screwPositionTs := make([]float64, hertz*5)
	injectScrewPosition := int(math.Ceil(float64(rel * 50)))
	manu1ScrewPosition := injectScrewPosition * 5
	manu2ScrewPosition := int(math.Ceil(float64(rel * 1600)))
	returnScrewPosition := int(math.Ceil(float64(rel * 1700)))

	machineInjectionTs := make([]float64, hertz*5)
	startInjection := int(math.Ceil(float64(rel * 100)))
	maxStep1Injection := int(math.Ceil(float64(rel * 120)))
	maxStep2Injection := int(math.Ceil(float64(rel * 200)))
	repressInjection := int(math.Ceil(float64(rel * 1600)))
	down1Pressure := int(math.Ceil(float64(rel * 1750)))
	down2Pressure := int(math.Ceil(float64(rel * 1800)))

	machinePressuseTs := make([]float64, hertz*5)
	startMaxPressure := int(math.Ceil(float64(rel * 200)))
	endMaxPressure := int(math.Ceil(float64(rel * 250)))

	machineTempTs := make([]float64, hertz*5)
	maxTemp := int(math.Ceil(float64(rel * 200)))

	return &PLCRunner{
		tWindow:             tWindow,
		screwPositionTs:     screwPositionTs,
		injectScrewPosition: injectScrewPosition,
		manu1ScrewPosition:  manu1ScrewPosition,
		manu2ScrewPosition:  manu2ScrewPosition,
		returnScrewPosition: returnScrewPosition,
		machineInjectionTs:  machineInjectionTs,
		startInjection:      startInjection,
		maxStep1Injection:   maxStep1Injection,
		maxStep2Injection:   maxStep2Injection,
		repressInjection:    repressInjection,
		down1Pressure:       down1Pressure,
		down2Pressure:       down2Pressure,
		machinePressuseTs:   machinePressuseTs,
		startMaxPressure:    startMaxPressure,
		endMaxPressure:      endMaxPressure,
		machineTempTs:       machineTempTs,
		maxTemp:             maxTemp,
		stopCh:              make(chan bool),
	}
}

func (p *PLCRunner) Stop() {
	p.stopCh <- true
}

func (p *PLCRunner) Run(selfClient OPCUAClient, lineClient OPCUAClient, log logr.Logger) {
	var i int = 0
	var oldMchOn bool = false

	go p.genValues(log)
	for {
		time.Sleep(time.Duration(p.tWindow) * time.Millisecond)
		select {
		case <-p.stopCh:
			log.Info("stopped plc")
			return
		default:
			mchOn, err := lineClient.ReadVar(OPC_NS_EXTERNAL, OPC_PATH_EXTERNAL, "f_mch_progress_1")
			if err != nil {
				log.Error(err, "Failed to read variable f_mch_progress_1 from line plc")
				continue
			}

			var errWrite error
			var vars map[string]interface{}
			if mchOn != nil && mchOn.(bool) {
				log.Info("Status=On")
				if i < len(p.screwPositionTs) {
					vars = p.getValues(i)
					i++
				} else {
					vars = p.getValuesZero(i)
				}
				errWrite = selfClient.WriteVar(OPC_NS_INTERNAL, "", vars)
			} else if oldMchOn {
				log.Info("Status=Paused")
				go p.genValues(log)
				vars = p.getValuesZero(i)
				errWrite = selfClient.WriteVar(OPC_NS_INTERNAL, "", vars)
				i = 0
			}

			if errWrite != nil {
				log.Error(errWrite, "Failed to write variable")
			}
			oldMchOn = mchOn.(bool)
		}
	}
}

func (p *PLCRunner) getValues(i int) map[string]interface{} {
	vars := make(map[string]interface{})
	vars["screwPosition"] = p.screwPositionTs[i]
	vars["injectPressure"] = p.machineInjectionTs[i]
	vars["pressure"] = p.machinePressuseTs[i]
	vars["temp"] = p.machineTempTs[i]
	return vars
}

func (p *PLCRunner) getValuesZero(i int) map[string]interface{} {
	vars := make(map[string]interface{})
	vars["screwPosition"] = float64(0)
	vars["injectPressure"] = float64(0)
	vars["pressure"] = float64(0)
	vars["temp"] = float64(0)
	return vars
}

func linspace(start, end float64, num int) []float64 {
	result := make([]float64, num)
	if num == 1 {
		result[0] = start
		return result
	}
	step := (end - start) / float64(num-1)
	for i := range result {
		result[i] = start + float64(i)*step
	}
	return result
}

func (p *PLCRunner) genValues(log logr.Logger) {
	log.Info("Generating values for next tWindow...")
	// Injection Screw Position
	copy(p.screwPositionTs[0:p.injectScrewPosition], linspace(2.8, -0.7, p.injectScrewPosition))
	// Manutenção de posição
	for i := p.injectScrewPosition; i < p.manu1ScrewPosition; i++ {
		p.screwPositionTs[i] = -0.7
	}
	for i := p.manu1ScrewPosition; i < p.manu2ScrewPosition; i++ {
		p.screwPositionTs[i] = -0.8
	}
	// Pequena variação (simulando recalque)
	for i := p.manu2ScrewPosition; i < p.returnScrewPosition; i++ {
		p.screwPositionTs[i] = -0.7
	}
	// Dosagem (retorno linear lento)
	copy(
		p.screwPositionTs[p.returnScrewPosition:],
		linspace(-0.7, 2.5, len(p.screwPositionTs)-p.returnScrewPosition),
	)

	// Pressão de Injeção
	for i := range p.machineInjectionTs {
		p.machineInjectionTs[i] = -0.6
	}
	// Pico de injeção
	for i := p.startInjection; i < p.maxStep1Injection; i++ {
		p.machineInjectionTs[i] = 4.5
	}
	for i := p.maxStep1Injection; i < p.maxStep2Injection; i++ {
		p.machineInjectionTs[i] = 0.5
	}
	// Pressão de recalque
	for i := p.maxStep2Injection; i < p.repressInjection; i++ {
		p.machineInjectionTs[i] = 0.4
	}
	// Queda de pressão
	for i := p.repressInjection; i < p.down1Pressure; i++ {
		p.machineInjectionTs[i] = -1.0
	}
	for i := p.down1Pressure; i < p.down2Pressure; i++ {
		p.machineInjectionTs[i] = -0.8
	}

	// Gráfico 3: Pressão no Molde
	for i := range p.machinePressuseTs {
		p.machinePressuseTs[i] = -2.0
	}
	// Preenchimento rápido do molde
	copy(
		p.machinePressuseTs[p.startMaxPressure:p.endMaxPressure],
		linspace(-2.0, 9.5, p.endMaxPressure-p.startMaxPressure),
	)
	// Queda exponencial da pressão
	decaimentoSpace := linspace(0, 5, len(p.machinePressuseTs)-p.endMaxPressure)
	for i := 0; i < len(p.machinePressuseTs)-p.endMaxPressure; i++ {
		p.machinePressuseTs[p.endMaxPressure+i] = 9.5*math.Exp(-decaimentoSpace[i]) - 2.0
	}
	// Adicionando ruído
	for i := range p.machinePressuseTs {
		p.machinePressuseTs[i] += rand.NormFloat64() * 0.05
	}

	// Subida rápida da temperatura
	subidaSpace := linspace(0, 10, p.maxTemp)
	for i := 0; i < p.maxTemp; i++ {
		p.machineTempTs[i] = 1.4*(1-math.Exp(-subidaSpace[i])) - 2.8
	}
	// Resfriamento lento e linear
	copy(p.machineTempTs[p.maxTemp:], linspace(1.4, -2.5, len(p.machineTempTs)-p.maxTemp))
	// Adicionando ruído significativo
	for i := range p.machineTempTs {
		p.machineTempTs[i] += rand.NormFloat64() * 0.15
	}
}
