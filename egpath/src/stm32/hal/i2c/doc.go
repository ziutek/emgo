// Package i2c provides interface to use I2C peripheral.
//
// I2C peripheral can be used through driver. Only one driver can be used with
// one peripheral.
//
// In case of operation as master device, driver implements virtual connections
// to slave devices.
//
// Driver supports multiple concurrent master connections. First read or write
// on inactive connection starts an I2C transaction and the connection becomes
// active until the transaction end. Peripheral supports only one active
// connection at the same time. Starting a subsequent transaction in other
// connection is blocked until the current transaction will end.
//
// Active connection supports both read and write transactions. There is no
// need to terminate write transaction before subsequent read transaction but
// read transaction must be terminated before subsequent write transaction.
package i2c