package net.getMeStore.client.service;

import com.fasterxml.jackson.core.JsonProcessingException;
import com.fasterxml.jackson.databind.ObjectMapper;
import net.getMeStore.client.config.Env;
import net.getMeStore.client.config.UdsHandler;
import net.getMeStore.client.models.AppendRequestPayload;
import org.springframework.http.MediaType;
import org.springframework.stereotype.Service;
import org.springframework.web.reactive.function.client.WebClient;
import reactor.core.publisher.Mono;

@Service
public class GetMeService {

    WebClient udsWebClient;
    Env env;

    ObjectMapper objectMapper;

    public GetMeService(UdsHandler udsHandler, Env env) {
        this.env = env;
        this.objectMapper = new ObjectMapper();
        this.udsWebClient = udsHandler.webClient(env.getSocketPath());
    }

    public Mono<String> get(String key) {
        return udsWebClient
                .get()
                .uri(uriBuilder -> uriBuilder.path("/get").queryParam("key", key).build())
                .retrieve()
                .bodyToMono(String.class);
    }

    public Mono<String> put(String key, String value) throws JsonProcessingException {
        return udsWebClient
                .post()
                .uri(uriBuilder -> uriBuilder.path("/put").build())
                .contentType(MediaType.APPLICATION_JSON)
                .bodyValue(objectMapper.writeValueAsString(new AppendRequestPayload(key, value)))
                .retrieve()
                .bodyToMono(String.class);
    }

    public Mono<String> delete (String key) {
        return udsWebClient
                .delete()
                .uri(uriBuilder -> uriBuilder.path("/delete").queryParam("key", key).build())
                .retrieve()
                .bodyToMono(String.class);
    }

    public Mono<String> clearStore() {
        return udsWebClient
                .delete()
                .uri(
                uriBuilder -> uriBuilder.path("/clear").build())
                .retrieve()
                .bodyToMono(String.class);
    }
}
