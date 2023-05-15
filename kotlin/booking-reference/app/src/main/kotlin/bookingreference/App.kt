package bookingreference

import io.vertx.core.Future
import io.vertx.core.Vertx
import io.vertx.core.Vertx.vertx
import io.vertx.core.json.JsonObject
import io.vertx.ext.web.Router
import java.lang.System.Logger.Level.*
import java.lang.System.getLogger
import java.lang.System.getenv
import java.util.concurrent.atomic.AtomicInteger

private const val defaultStartingPoint = 12345678
private val logger = getLogger("main")

fun main() {
    startApplication(8080)
}

internal fun startApplication(httpPort: Int): Vertx {
    val vertx = vertx()
    val referenceCounter = ReferenceCounter(getReferenceStartingPoint())

    val router = Router.router(vertx)
    router.get("/").respond {
        Future.succeededFuture(JsonObject.of("value", referenceCounter.getNextReference()))
    }

    logger.log(INFO, "Starting booking reference service on port {0,number,#}.", httpPort)
    vertx.createHttpServer()
        .requestHandler(router)
        .listen(httpPort)
        .onFailure {
            logger.log(ERROR, "An error occurred starting the service. Cause: ${it.message}.")
            vertx.close()
        }
        .onComplete {
            logger.log(INFO, "Closing service. Bye.")
        }

    return vertx
}

private class ReferenceCounter(startingPoint: Int) {
    private val counter = AtomicInteger(startingPoint)
    fun getNextReference(): String {
        val current = counter.incrementAndGet()
        return String.format("%016x", current)
    }
}

fun getReferenceStartingPoint(): Int {
    val startingPoint = getenv("STARTING_POINT")?.toIntOrNull()
    if (startingPoint == null) {
        logger.log(WARNING, "Unable to obtain reference starting point from STARTING_POINT using default.")
    }
    return startingPoint ?: defaultStartingPoint
}
