package net.getMeStore.client.controllers;

import com.fasterxml.jackson.core.JsonProcessingException;
import net.getMeStore.client.service.GetMeService;
import org.springframework.web.bind.annotation.*;
import reactor.core.publisher.Mono;


import net.getMeStore.client.models.BatchGetResult;
import net.getMeStore.client.models.BatchPutResult;
import net.getMeStore.client.models.BatchDeleteResult;

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

    @GetMapping("/batch-get")
    public Mono<BatchGetResult> batchGet(@RequestParam String jsonPath) {
        return this.getMeService.batchGet(jsonPath);
    }

    @PostMapping("/put")
    public Mono<String> put(@RequestParam String key, @RequestParam String value) throws JsonProcessingException {
        return this.getMeService.put(key, value);
    }

    @PostMapping("/batch-put")
    public Mono<BatchPutResult> batchPut(@RequestBody String jsonPayload) {
        return this.getMeService.batchPut(jsonPayload);
    }

    @DeleteMapping("/delete")
    public Mono<String> delete(@RequestParam String key) {
        return this.getMeService.delete(key);
    }

    @DeleteMapping("/batch-delete")
    public Mono<BatchDeleteResult> batchDelete(@RequestBody String jsonPayload) {
        return this.getMeService.batchDelete(jsonPayload);
    }

    @DeleteMapping("/clearStore")
    public Mono<String> clearStore() {
        return this.getMeService.clearStore();
    }

}
