package connpool

type idleList struct {
	count       int
	front, back IConn
}

func (l *idleList) pushFront(pc IConn) {
	pc.setNext(l.front)
	pc.setPrev(nil)
	if l.count == 0 {
		l.back = pc
	} else {
		l.front.setPrev(pc)
	}
	l.front = pc
	l.count++
}

func (l *idleList) popFront() {
	pc := l.front
	l.count--
	if l.count == 0 {
		l.front, l.back = nil, nil
	} else {
		pc.getNext().setPrev(nil)
		l.front = pc.getNext()
	}
	pc.setNext(nil)
	pc.setPrev(nil)
}

func (l *idleList) popBack() {
	pc := l.back
	l.count--
	if l.count == 0 {
		l.front, l.back = nil, nil
	} else {
		pc.getPrev().setNext(nil)
		l.back = pc.getPrev()
	}
	pc.setNext(nil)
	pc.setPrev(nil)
}
