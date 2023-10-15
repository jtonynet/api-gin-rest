import io.gatling.core.Predef._
import io.gatling.http.Predef._
import scala.concurrent.duration._

class Alunos extends Simulation {

  val httpProtocol = http
    .baseUrl("http://nginx-gin-rest:80")
    .acceptHeader("text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
    .acceptLanguageHeader("en-US,en;q=0.5")
    .acceptEncodingHeader("gzip, deflate")
    .userAgentHeader("Mozilla/5.0 (Macintosh; Intel Mac OS X 10.8; rv:16.0) Gecko/20100101 Firefox/16.0")


  val alunos = 
    feed(tsv("alunos.tsv").circular())
    .exec(session => {
      val payload = session("payload").as[String]
      println(s"Payload JSON: $payload")
      session
    })
    .exec(http("Requisição para /alunos")
      .post("/aluno")
      .body(StringBody("#{payload}"))
      .header("Content-Type", "application/json")
      .check(status.in(200,202))
    )


  val testAlunos = scenario("Teste Alunos").exec(alunos)

  setUp(
    testAlunos.inject(rampUsers(10000).during(10.seconds))
  ).protocols(httpProtocol)
}
