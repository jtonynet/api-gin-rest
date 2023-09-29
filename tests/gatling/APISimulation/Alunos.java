package APISimulation;

import static io.gatling.javaapi.core.CoreDsl.*;
import static io.gatling.javaapi.http.HttpDsl.*;

import io.gatling.javaapi.core.*;
import io.gatling.javaapi.http.*;

public class Alunos extends Simulation {

    HttpProtocolBuilder httpProtocol =
        http.baseUrl("http://nginx-gin-rest:80")
            .acceptHeader("text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
            .acceptLanguageHeader("en-US,en;q=0.5")
            .acceptEncodingHeader("gzip, deflate")
            .userAgentHeader(
                "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.8; rv:16.0) Gecko/20100101 Firefox/16.0"
            );

    ChainBuilder alunos =
        exec(http("Requisição para /alunos")
            .get("/alunos")
            .check(status().is(200))
        );

    ScenarioBuilder testAlunos = scenario("Teste Alunos").exec(alunos);

    {
        setUp(
            testAlunos.injectOpen(rampUsers(10).during(10))
        ).protocols(httpProtocol);
    }
}
