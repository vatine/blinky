import protos.blinky_pb2_grpc as blinky_grpc
import protos.blinky_pb2 as blinky

import blinkstick.blinkstick as blinkstick

import concurrent.futures


class BlinkstickNotFound(exception):
    def __str__(self):
        return "No blinkstick found."


class BlinkyServer(blinky_grpc.BlinkyServicer):
    def __init__(self):
        super().__init__(self)

        sticks = blinkstick.find_all()
        if not sticks:
            raise BlinkstickNotFound()

        self._stick = sticks[0]

    def SetLEDs(self, request, context):
        leds = request.LEDs
        response = blinky.SetLEDResponse(LEDs=leds)

        red = clamp(request.Red)
        blue = clamp(request.Blue)
        green = clamp(request.Green)

        for led in leds:
            self._stick.set_Color(index=led, red=red, blue=blue, green=green)

        return response

    def GetLEDs(self, request, context):
        leds = request.LEDs
        if not leds:
            leds = range(self._stick.get_led_count())

        statuses =[]

        for led in leds:
            red, green, blue = self._stick.get_color(index=led)
            statuses.append(blinky.LEDStatus(LED=led, Red=red, Green=green, Blue=blue))

        return blinky.GetLEDResponse(Status=statuses)


def clamp(number):
    """Return a number clamped to between 0 and 255."""

    if number <= 0:
        return 0
    if number >= 255:
        return 255

    return number

def serve():
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    blinky_grpc.add_BlinkyServicer_to_server(BlinkyServer(), server)
    server.add_insecure_port(":4004")
    server.start()
    server.wait_for_termination()


if __name__ == '__main__':
    serve()
