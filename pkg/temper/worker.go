/**
 * Copyright 2021 SAP SE
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package temper

import (
	"sync"

	"github.com/sapcc/baremetal_temper/pkg/node"
)

type JobChannel chan *node.Node
type JobQueue chan chan *node.Node

type Worker struct {
	JobChan JobChannel
	Queue   JobQueue
	Quit    chan struct{}
}

func (w *Worker) Start() {
	go func() {
		for {
			w.Queue <- w.JobChan
			select {
			case job := <-w.JobChan:
				var wg sync.WaitGroup
				job.AddAllTemperTasks(true, true, true, true)
				wg.Add(1)
				job.Temper(true, &wg)
				wg.Wait()
			case <-w.Quit:
				close(w.JobChan)
				return
			}
		}
	}()
}