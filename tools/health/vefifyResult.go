package main

type VerifyResult uint8

const (
	/** Peer was verified. */
	success = uint8(0)

	/** An i/o error was encountered during verification. */
	ioError = uint8(1)

	/** Peer sent malformed data. */
	malformedData = uint8(2)

	/** Peer failed the challenge. */
	failedChallenge = uint8(3)
)
