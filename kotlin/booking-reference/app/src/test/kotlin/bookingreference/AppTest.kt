package bookingreference

import io.vertx.core.json.JsonObject
import io.vertx.ext.web.client.WebClient
import java.net.ServerSocket
import kotlin.test.Test
import kotlin.test.assertEquals

class AppTest {

    @Test
    fun startsService() {
        val httpPort = getFreePort()
        val vertx = startApplication(httpPort)

        val client = WebClient.create(vertx)
        val result = client.get(httpPort, "localhost", "/")
            .send()
            // FIXME: Use vertx unit testing framework
            .toCompletionStage()
            .toCompletableFuture()
            .join()

        assertEquals(200, result.statusCode())
        assertEquals(JsonObject.of("value", "0000000000bc614f"), result.bodyAsJsonObject())

        vertx.close()
    }

    private fun getFreePort(): Int {
        ServerSocket(0).use { socket ->
            return socket.localPort
        }
    }
}
