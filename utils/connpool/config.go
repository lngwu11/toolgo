package connpool

type Config struct {
	// Factory is an application supplied function for creating and configuring a connection.
	Factory func(p interface{}) (IConn, error)
	// Param is the parameter of Factory function
	Param interface{}

	// Maximum number of idle connections in the pool.
	MaxIdle int

	// Maximum number of connections allocated by the pool at a given time.
	// When zero, there is no limit on the number of connections in the pool.
	MaxActive int

	// If Wait is true and the pool is at the MaxActive limit, then Get() waits
	// for a connection to be returned to the pool before returning.
	Wait bool
}
