package hosting

// IHostApplicationLifetime provides notifications for application lifetime events.
type IHostApplicationLifetime interface {
	// ApplicationStarted returns a channel that is closed when the application has fully started.
	ApplicationStarted() <-chan struct{}

	// ApplicationStopping returns a channel that is closed when the application is stopping.
	ApplicationStopping() <-chan struct{}

	// ApplicationStopped returns a channel that is closed when the application has stopped.
	ApplicationStopped() <-chan struct{}

	// StopApplication requests the application to stop.
	StopApplication()

	// Internal notification methods
	NotifyStarting()
	NotifyStarted()
	NotifyStopping()
	NotifyStopped()
}

// ApplicationLifetime implements IHostApplicationLifetime.
type ApplicationLifetime struct {
	startingChan chan struct{}
	startedChan  chan struct{}
	stoppingChan chan struct{}
	stoppedChan  chan struct{}
}

// NewApplicationLifetime creates a new ApplicationLifetime.
func NewApplicationLifetime() IHostApplicationLifetime {
	return &ApplicationLifetime{
		startingChan: make(chan struct{}),
		startedChan:  make(chan struct{}),
		stoppingChan: make(chan struct{}),
		stoppedChan:  make(chan struct{}),
	}
}

func (l *ApplicationLifetime) ApplicationStarted() <-chan struct{} {
	return l.startedChan
}

func (l *ApplicationLifetime) ApplicationStopping() <-chan struct{} {
	return l.stoppingChan
}

func (l *ApplicationLifetime) ApplicationStopped() <-chan struct{} {
	return l.stoppedChan
}

func (l *ApplicationLifetime) StopApplication() {
	select {
	case <-l.stoppingChan:
		// Already closed
	default:
		close(l.stoppingChan)
	}
}

func (l *ApplicationLifetime) NotifyStarting() {
	// Can trigger events here if needed
}

func (l *ApplicationLifetime) NotifyStarted() {
	close(l.startedChan)
}

func (l *ApplicationLifetime) NotifyStopping() {
	select {
	case <-l.stoppingChan:
		// Already closed
	default:
		close(l.stoppingChan)
	}
}

func (l *ApplicationLifetime) NotifyStopped() {
	close(l.stoppedChan)
}

