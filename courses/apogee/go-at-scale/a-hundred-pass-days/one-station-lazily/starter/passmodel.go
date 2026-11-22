package main

// passModel is the pass model the planner asks for each day's passes.
// Tests put their own in its place, to watch, hold or cancel a run.
var passModel Predictor = predict
