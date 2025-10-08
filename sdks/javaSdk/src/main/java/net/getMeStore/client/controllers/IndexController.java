package net.getMeStore.client.controllers;

import com.fasterxml.jackson.core.JsonProcessingException;
import net.getMeStore.client.service.GetMeService;
import org.springframework.web.bind.annotation.*;
import reactor.core.publisher.Mono;

@RestController
public class IndexController {

    private final GetMeService getMeService;

    IndexController(GetMeService getMeService) {
        this.getMeService = getMeService;
    }

    @GetMapping("/get")
    public Mono<String> get(@RequestParam String key) {
        return this.getMeService.get(key);
    }

    @PostMapping("/put")
    public Mono<String> put(@RequestParam String key, @RequestParam String value) throws JsonProcessingException {
        return this.getMeService.put(key, value);
    }

    @DeleteMapping("/delete")
    public Mono<String> delete(@RequestParam String key) {
        return this.getMeService.delete(key);
    }

    @DeleteMapping("/clearStore")
    public Mono<String> clearStore() {
        return this.getMeService.clearStore();
    }

}
