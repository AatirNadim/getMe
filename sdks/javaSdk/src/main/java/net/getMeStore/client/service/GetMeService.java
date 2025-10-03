package net.getMeStore.client.service;

import com.fasterxml.jackson.core.JsonProcessingException;
import com.fasterxml.jackson.databind.ObjectMapper;
import net.getMeStore.client.config.Env;
import net.getMeStore.client.config.UdsHandler;
import net.getMeStore.client.models.AppendRequestPayload;
import org.springframework.http.MediaType;
import org.springframework.stereotype.Service;
import reactor.core.publisher.Mono;

@Service
public class GetMeService {

    UdsHandler udsHandler;
    Env env;

    ObjectMapper objectMapper;

    public GetMeService(UdsHandler udsHandler, Env env) {
        this.udsHandler = udsHandler;
        this.env = env;
        this.objectMapper = new ObjectMapper();
    }

    public Mono<String> get(String key) {
        return udsHandler.webClient(env.getSocketPath())
                .get()
                .uri(uriBuilder -> uriBuilder.path("/get").queryParam("key", key).build())
                .retrieve()
                .bodyToMono(String.class);
    }

    public Mono<String> put(String key, String value) throws JsonProcessingException {


        return udsHandler.webClient(env.getSocketPath())
                .post()
                .uri(uriBuilder -> uriBuilder.path("/put").build())
                .contentType(MediaType.APPLICATION_JSON)
                .bodyValue(objectMapper.writeValueAsString(new AppendRequestPayload(key, value)))
                .retrieve()
                .bodyToMono(String.class);
    }

    public Mono<String> delete (String key) {
        return udsHandler.webClient(env.getSocketPath())
                .delete()
                .uri(uriBuilder -> uriBuilder.path("/delete").queryParam("key", key).build())
                .retrieve()
                .bodyToMono(String.class);
    }
}
