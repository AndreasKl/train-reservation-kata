package bookingreference

import io.vertx.core.Vertx;

fun main() {
    val server = Vertx.vertx().createHttpServer()
    server.listen(8080, "0.0.0.0")
    
}
