package main

type WaterValues struct {
	targetHeight float64
	height       float64
	speed        float64
}

func (waterValues *WaterValues) Construct(height float64, targetHeight float64, speed float64) {
	waterValues.height = height
	waterValues.targetHeight = targetHeight
	waterValues.speed = speed
}

func (waterValues *WaterValues) Update(dampening float64, tension float64) {
	x := waterValues.targetHeight - waterValues.height
	waterValues.speed += tension*x - waterValues.speed*dampening
	waterValues.height += waterValues.speed
}
