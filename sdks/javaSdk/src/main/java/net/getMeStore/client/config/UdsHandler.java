package net.getMeStore.client.config;

import io.netty.channel.unix.DomainSocketAddress;
import org.springframework.context.annotation.Configuration;
import org.springframework.http.client.reactive.ReactorClientHttpConnector;
import org.springframework.web.reactive.function.client.WebClient;
import reactor.netty.http.client.HttpClient;

@Configuration
public class UdsHandler {

    public WebClient webClient(String socketPath) {
        DomainSocketAddress domainSocketAddress = new DomainSocketAddress(socketPath);

        HttpClient httpClient = HttpClient.create().remoteAddress(() -> domainSocketAddress);

        ReactorClientHttpConnector connector = new ReactorClientHttpConnector(httpClient);

        return WebClient.builder().clientConnector(connector).build();
    }
}
